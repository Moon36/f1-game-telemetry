package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/moon36/f1-game-telemetry/src/internal/packets"
	v23 "github.com/moon36/f1-game-telemetry/src/internal/packets/v23"
	"github.com/redis/go-redis/v9"
)

/*
Converts a byte array to a string by trimming any trailing null bytes.

Parameters:
  - b: The byte array to be converted to a string.
*/
func byteArrayToString(b []byte) string {
	return string(bytes.TrimRight(b, "\x00"))
}

/*
Looks for the given ID in the given hash-set in Redis and returns the associated value.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - hashSet: The hash-set/key to use for the lookup.
  - id: The field to look up.
*/
func getNameById(ctx context.Context, redisClient *redis.Client, hashSet string, id uint8) (string, error) {
	value, err := redisClient.HGet(ctx, hashSet, fmt.Sprintf("%d", id)).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

/*
Resolves the IDs in the given ParticipantData struct to their corresponding names using the data stored in Redis.
The enriched participant data is then sent through the provided channel along with any potential error that occurred.

Parameters:
  - ctx: The context for managing the lifecycle of Redis operations.
  - redisClient: The Redis client used to interact with the Redis database.
  - participantId: The ID of the participant being enriched, used for the channel data struct.
  - participant: The ParticipantData struct containing the participant data with IDs to be resolved.
  - ch: The channel through which the enriched participant data and any potential error will be sent.
*/
func enrichParticipant23(ctx context.Context,
	redisClient *redis.Client,
	participantId uint8,
	participant v23.ParticipantData,
	ch chan ChannelData) {

	storeParticipant := packets.StoreParticipant{
		M_aiControlled:    participant.M_aiControlled,
		M_networkId:       participant.M_networkId,
		M_myTeam:          participant.M_myTeam,
		M_raceNumber:      participant.M_raceNumber,
		M_name:            byteArrayToString(participant.M_name[:]),
		M_yourTelemetry:   participant.M_yourTelemetry,
		M_showOnlineNames: participant.M_showOnlineNames,
	}
	err := error(nil)
	// Driver ID to name
	driverName := "Player"
	if participant.M_driverId != 255 {
		driverName, err = getNameById(ctx, redisClient, "csv:drivers", participant.M_driverId)
		if err != nil {
			ch <- ChannelData{
				Idx:             participantId,
				ParticipantData: packets.StoreParticipant{},
				Error:           err,
			}
			return
		}
	}
	storeParticipant.M_driverName = driverName

	// Team ID to name
	teamName := ""
	if participant.M_teamId != 255 {
		teamName, err = getNameById(ctx, redisClient, "csv:teams", participant.M_teamId)
		if err != nil {
			ch <- ChannelData{
				Idx:             participantId,
				ParticipantData: packets.StoreParticipant{},
				Error:           err,
			}
			return
		}
	}
	storeParticipant.M_teamName = teamName

	// Nationality ID to name
	nationality := ""
	if participant.M_nationality != 255 {
		nationality, err = getNameById(ctx, redisClient, "csv:nationalities", participant.M_nationality)
		if err != nil {
			ch <- ChannelData{
				Idx:             participantId,
				ParticipantData: packets.StoreParticipant{},
				Error:           err,
			}
			return
		}
	}
	storeParticipant.M_nationality = nationality

	// Platform ID to name
	platformName := ""
	if participant.M_platform != 0 {
		platformName, err = getNameById(ctx, redisClient, "csv:platforms", participant.M_platform)
		if err != nil {
			ch <- ChannelData{
				Idx:             participantId,
				ParticipantData: packets.StoreParticipant{},
				Error:           err,
			}
			return
		}
	}
	storeParticipant.M_platformName = platformName

	ch <- ChannelData{
		Idx:             participantId,
		ParticipantData: storeParticipant,
		Error:           nil,
	}
}
