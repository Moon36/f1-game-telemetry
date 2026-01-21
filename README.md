# F1 Game Telemetry Dashboard
A dashboard for visualizing telemetry data from F1 games.

## Currently supported Games
- F1 2023


## AI Race Engineer
Of course we have to integrate AI into this... otherwise it wouldn't be a modern software project, right?\
The AI Race Engineer is a **local** LLM model which is running in the background, which will receive some of your data to give you (hopefully) useful information.\

> It is more thought to be a fun gimmick than an actual race engineer with tactical and statistical expertise.\
> Think of it like being assisted by a [Ferrari race engineer](https://youtu.be/vi0vMUFEZrM?si=Bp3ebqmXVb21PuCb)

### Using the AI Race Engineer
For the race engineer to work, make sure you have a Docker Version, which supports the **Docker Model Runner** feature.

To install a model use:
```shell
$ docker model pull <MODEL_NAME>
```
Where `<MODEL_NAME>` is the model you want to install. For computers or especially laptops with minimal specs, I'd suggest the `ai/smollm2` model. It is very small (about 256MB) and *should* run even on relatively low-spec systems.

You can list your installed LLM models via:
```shell
$ docker model list
MODEL NAME          PARAMETERS  QUANTIZATION    ARCHITECTURE   MODEL ID      CREATED        CONTEXT  SIZE
granite-4.0-h-nano  1.46 B      MOSTLY_Q8_0     granitehybrid  91eb206d6605  2 months ago            1.45 GiB
smollm2             361.82 M    IQ2_XXS/Q4_K_M  llama          354bf30d0aa3  10 months ago           256.35 MiB
```

Then just set the `LLM_MODEL` in the `.env` file to the model you want to use.
For example:
```
LLM_MODEL=ai/smollm2
```

> Keep in mind though, that since the `ai/smollm2` model is so small, that its output might not be the most "analytically correct".\
> Think of it as a toddler relaying some statistical information to you 😄


## Current Development Status
- [X] Implement UDP server to receive telemetry data
- [X] Integrate message broker
- [X] Implement message consumer
- [ ] Add statistics collector
- [X] Add frontend backend
- [X] Develop frontend dashboard
