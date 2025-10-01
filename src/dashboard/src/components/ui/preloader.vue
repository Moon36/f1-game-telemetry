<!-- original: https://codepen.io/kol123/pen/ALgEmQ -->

<template>
  <div class="preloader-overlay" v-if="isVisible" :class="{ 'exiting': isExiting }">
    <div class="preloader-content" :class="{ 'exiting': isExiting }">
      <div id="Container" :class="{ 'exiting': isExiting }">
        <div id="nose-top"></div>
        <div id="nose-bottom"></div>
        <div id="nose"></div>
        <div id="front-wing"></div>
        <div id="top-front-wing-trim"></div>
        <div id="bottom-front-wing-trim"></div>
        <div id="top-front-wing-trim-2"></div>
        <div id="bottom-front-wing-trim-2"></div>
        <div id="top-front-wing"></div>
        <div id="top-front-wing-tail"></div>
        <div id="bottom-front-wing"></div>
        <div id="bottom-front-wing-tail"></div>
        <div id="bottom-front-wheel"></div>
        <div id="bottom-back-wheel"></div>
        <div id="top-front-wheel"></div>
        <div id="top-back-wheel"></div>
        <div id="rear-body"></div>
        <div id="rear-wing-bg"></div>
        <div id="rear-wing"></div>
        <div id="top-body-curve"></div>
        <div id="top-body-curve-cut"></div>
        <div id="top-body-curve-straight"></div>
        <div id="top-body-curve-straight-2"></div>
        <div id="bottom-body-curve"></div>
        <div id="bottom-body-curve-cut"></div>
        <div id="bottom-body-curve-straight"></div>
        <div id="bottom-body-curve-straight-2"></div>
        <div id="back-body-curve"></div>
        <div id="body-hood"></div>
        <div id="back-body"></div>
        <div id="back-body-top"></div>
        <div id="back-body-bottom"></div>
        <div id="back-body-2"></div>
        <div id="top-spoke-1"></div>
        <div id="top-spoke-2"></div>
        <div id="top-spoke-3"></div>
        <div id="top-spoke-4"></div>
        <div id="bottom-spoke-1"></div>
        <div id="bottom-spoke-2"></div>
        <div id="bottom-spoke-3"></div>
        <div id="bottom-spoke-4"></div>
        <div id="back-spoke"></div>
        <div id="mirror-top"></div>
        <div id="mirror-bottom"></div>
        <div id="driver-bg"></div>
        <div id="driver-wheel"></div>
        <div id="driver-helmet"></div>
        <div id="bottom-body-spine"></div>
        <div id="top-body-spine"></div>
        <div id="end-body-spine"></div>
        <div id="top-body-spine-2"></div>
        <div id="bottom-body-spine-2"></div>
      </div>
      
      <div class="loading-text" :class="{ 'exiting': isExiting }">
        <!--<h2>F1 Telemetry Dashboard</h2>-->
        <p>{{ loadingMessage }}</p>
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
        <div class="progress-text">{{ Math.round(progress) }}%</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Props {
  duration?: number
  autoHide?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  duration: 3000,
  autoHide: true
})

const emit = defineEmits<{
  complete: []
}>()

const isVisible = ref(true)
const isExiting = ref(false)
const progress = ref(0)
const loadingMessage = ref('Initializing Dashboard...')

const messages = [
  'Initializing Dashboard...',
  'Connecting to Telemetry...',
  'Loading Car Data...',
  'Preparing Interface...',
  'Almost Ready...'
]

let progressInterval: number | undefined
let messageInterval: number | undefined

const startLoading = () => {
  let messageIndex = 0
  
  // Progress animation
  progressInterval = window.setInterval(() => {
    if (progress.value < 100) {
      progress.value += Math.random() * 15 + 5 // Random increment between 5-20
      if (progress.value > 100) progress.value = 100
    }
  }, 200)

  // Message rotation
  messageInterval = window.setInterval(() => {
    messageIndex = (messageIndex + 1) % messages.length
    loadingMessage.value = messages[messageIndex]
  }, 800)

  // Auto-hide after duration
  if (props.autoHide) {
    setTimeout(() => {
      hidePreloader()
    }, props.duration)
  }
}

const hidePreloader = () => {
  isExiting.value = true
  // Start exit animation, then hide after animation completes
  setTimeout(() => {
    isVisible.value = false
    cleanup()
    emit('complete')
  }, 1500) // Wait longer for car to fully exit
}

const cleanup = () => {
  if (progressInterval) {
    clearInterval(progressInterval)
    progressInterval = undefined
  }
  if (messageInterval) {
    clearInterval(messageInterval)
    messageInterval = undefined
  }
}

// Public method to manually hide preloader
const hide = () => {
  hidePreloader()
}

onMounted(() => {
  startLoading()
})

onUnmounted(() => {
  cleanup()
})

// Expose hide method for parent components
defineExpose({
  hide
})
</script>

<script lang="ts">
export default {
  name: 'AppPreloader'
}
</script>

<style scoped>
.preloader-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  overflow: hidden;
  transition: opacity 0.5s ease-out 0.8s; /* Delay background fade */
}

.preloader-overlay.exiting {
  opacity: 0;
}

.preloader-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  transition: opacity 0.4s ease-out 0.6s; /* Delay content fade */
}

.preloader-content.exiting {
  opacity: 0;
}

#Container {
  width: 400px;
  height: 180px;
  position: relative;
  animation: entrance 2s ease-out, float 3s ease-in-out infinite 2s;
  transform-origin: center;
  transform: scale(0.4);
}

#Container.exiting {
  animation: exit 1.2s ease-in forwards;
}

#Container div {
  position: absolute;
  opacity: 1;
}

@keyframes entrance {
  0% { 
    transform: translateX(-150%) scale(0.4);
    opacity: 0;
  }
  100% { 
    transform: translateX(-50px) scale(0.4);
    opacity: 1;
  }
}

@keyframes exit {
  0% { 
    transform: translateX(-50px) scale(0.4);
    opacity: 1;
  }
  100% { 
    transform: translateX(150%) scale(0.4);
    opacity: 0;
  }
}

@keyframes float {
  0%, 100% { 
    transform: translateX(-50px) translateY(0px) scale(0.4);
  }
  50% { 
    transform: translateX(-50px) translateY(-6px) scale(0.42);
  }
}

/* Remove wheel animation - looks better static at small scale */

.loading-text {
  text-align: center;
  color: white;
  max-width: 250px;
  transition: opacity 0.4s ease-out 0.4s, transform 0.4s ease-out 0.4s; /* Delay text fade */
}

.loading-text.exiting {
  opacity: 0;
  transform: translateY(15px);
}

.loading-text h2 {
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.6rem;
  background: linear-gradient(45deg, #e74c3c, #f39c12);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.loading-text p {
  font-size: 0.85rem;
  margin-bottom: 1.2rem;
  opacity: 0.9;
  animation: fadeInOut 2s ease-in-out infinite;
}

@keyframes fadeInOut {
  0%, 100% { opacity: 0.7; }
  50% { opacity: 1; }
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #e74c3c, #f39c12);
  border-radius: 2px;
  transition: width 0.3s ease;
  position: relative;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.progress-text {
  font-size: 0.75rem;
  opacity: 0.8;
  font-weight: 500;
}

/* F1 Car Styles */
#nose {
  width: 400px;
  height: 25px;
  border-top-right-radius: 50px;
  border-bottom-right-radius: 50px;
  top: 156px;
  left: 190px;
  z-index: 2;
  background: #e74c3c;
}

#nose-top {
  width: 0;
  height: 0;
  border-bottom: 21px solid #e74c3c;
  border-left: 0px solid transparent;
  border-right: 200px solid transparent;
  top: 137px;
  left: 340px;
  z-index: 3;
}

#nose-bottom {
  width: 0;
  height: 0;
  border-top: 21px solid #e74c3c;
  border-left: 0px solid transparent;
  border-right: 200px solid transparent;
  top: 179px;
  left: 340px;
  z-index: 3;
}

#front-wing {
  width: 22px;
  height: 166px;
  top: 86px;
  left: 539px;
  z-index: 1;
  background: #e74c3c;
}

#top-front-wing {
  width: 0;
  height: 0;
  border-right: 30px solid #e74c3c;
  border-top: 0px solid transparent;
  border-bottom: 70px solid transparent;
  top: 86px;
  left: 515px;
  z-index: 1;
}

#top-front-wing-tail {
  width: 49px;
  height: 20px;
  background: #e74c3c;
  top: 86px;
  left: 512px;
}

#top-front-wheel {
  width: 60px;
  height: 35px;
  background: #222;
  top: 83px;
  left: 437px;
  border-radius: 8px;
  z-index: 2;
}

#top-back-wheel {
  width: 60px;
  height: 35px;
  background: #222;
  top: 83px;
  left: 120px;
  border-radius: 8px;
  z-index: 2;
}

#top-body-curve {
  width: 85px;
  height: 32px;
  background: #e74c3c;
  border-radius: 100px / 50px;
  top: 104px;
  left: 289px;
}

#top-body-curve-cut {
  width: 0;
  height: 0;
  border-left: 0px solid transparent;
  border-right: 18px solid transparent;
  border-bottom: 50px solid #e74c3c;
  top: 113px;
  left: 372px;
}

#top-body-curve-straight {
  width: 100px;
  height: 30px;
  background: #e74c3c;
  top: 116px;
  left: 207px;
  transform: rotate(-12deg);
  z-index: 2;
}

#top-body-curve-straight-2 {
  width: 45px;
  height: 30px;
  background: #e74c3c;
  top: 137px;
  left: 174px;
  transform: rotate(-36deg);
  z-index: 2;
}

#top-front-wing-trim {
  width: 14px;
  height: 6px;
  background: #e74c3c;
  top: 124px;
  left: 554px;
  z-index: 1;
}

#top-front-wing-trim-2 {
  width: 14px;
  height: 25px;
  background: #e74c3c;
  top: 92px;
  left: 554px;
  z-index: 1;
}

#bottom-front-wing-trim-2 {
  width: 14px;
  height: 25px;
  background: #e74c3c;
  top: 220px;
  left: 554px;
  z-index: 1;
}

#bottom-front-wing-trim {
  width: 14px;
  height: 6px;
  background: #e74c3c;
  top: 207px;
  left: 554px;
  z-index: 1;
}

#top-spoke-1 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 96px;
  left: 465px;
  transform: rotate(-9deg);
}

#top-spoke-2 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 105px;
  left: 475px;
  transform: rotate(-25deg);
}

#top-spoke-3 {
  width: 5px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 105px;
  left: 457px;
  transform: rotate(23deg);
}

#top-spoke-4 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 105px;
  left: 445px;
  transform: rotate(36deg);
}

#bottom-spoke-1 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 172px;
  left: 465px;
  transform: rotate(9deg);
}

#bottom-spoke-2 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 172px;
  left: 475px;
  transform: rotate(25deg);
}

#bottom-spoke-3 {
  width: 5px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 172px;
  left: 457px;
  transform: rotate(-23deg);
}

#bottom-spoke-4 {
  width: 8px;
  height: 60px;
  background: #777;
  z-index: 0;
  top: 172px;
  left: 445px;
  transform: rotate(-36deg);
}

#back-spoke {
  width: 18px;
  height: 160px;
  background: #777;
  z-index: 0;
  top: 92px;
  left: 141px;
}

#bottom-front-wing {
  width: 0;
  height: 0;
  border-right: 30px solid #e74c3c;
  border-bottom: 0px solid transparent;
  border-top: 70px solid transparent;
  top: 182px;
  left: 515px;
  z-index: 1;
}

#bottom-front-wing-tail {
  width: 49px;
  height: 20px;
  background: #e74c3c;
  top: 232px;
  left: 512px;
}

#bottom-front-wheel {
  width: 60px;
  height: 35px;
  background: #222;
  top: 219px;
  left: 437px;
  border-radius: 8px;
  z-index: 2;
}

#bottom-back-wheel {
  width: 60px;
  height: 35px;
  background: #222;
  top: 219px;
  left: 120px;
  border-radius: 8px;
  z-index: 2;
}

#rear-body {
  width: 16px;
  height: 96px;
  background: #e74c3c;
  top: 120px;
  left: 147px;
  z-index: 2;
}

#rear-wing-bg {
  width: 53px;
  height: 84px;
  border-top: 6px solid #e74c3c;
  border-bottom: 6px solid #e74c3c;
  background: #ecf0f1;
  top: 120px;
  left: 103px;
  z-index: 1;
}

#bottom-body-curve {
  width: 85px;
  height: 32px;
  background: #e74c3c;
  border-radius: 100px / 50px;
  top: 201px;
  left: 289px;
}

#bottom-body-curve-cut {
  width: 0;
  height: 0;
  border-left: 0px solid transparent;
  border-right: 18px solid transparent;
  border-top: 50px solid #e74c3c;
  top: 174px;
  left: 372px;
}

#bottom-body-curve-straight {
  width: 100px;
  height: 30px;
  background: #e74c3c;
  top: 191px;
  left: 207px;
  transform: rotate(12deg);
  z-index: 2;
}

#bottom-body-curve-straight-2 {
  width: 45px;
  height: 30px;
  background: #e74c3c;
  top: 171px;
  left: 174px;
  transform: rotate(36deg);
  z-index: 2;
}

#body-hood {
  width: 134px;
  height: 93px;
  background: #e74c3c;
  top: 123px;
  left: 240px;
  z-index: 2;
}

#back-body-curve {
  width: 85px;
  height: 60px;
  background: #e74c3c;
  border-radius: 100px / 50px;
  top: 139px;
  left: 168px;
  z-index: 1;
}

#back-body {
  width: 85px;
  height: 94px;
  background: #222;
  top: 122px;
  left: 148px;
  z-index: 0;
}

#back-body-top {
  width: 0;
  height: 0;
  border-left: 20px solid transparent;
  border-right: 0 solid transparent;
  border-bottom: 12px solid #222;
  z-index: 0;
  top: 111px;
  left: 172px;
}

#back-body-bottom {
  width: 0;
  height: 0;
  border-left: 20px solid transparent;
  border-right: 0 solid transparent;
  border-top: 12px solid #222;
  z-index: 0;
  top: 214px;
  left: 172px;
}

#back-body-2 {
  width: 97px;
  height: 115px;
  background: #222;
  top: 111px;
  left: 192px;
  z-index: 0;
}

#mirror-bottom {
  background: #e74c3c;
  width: 13px;
  height: 23px;
  top: 191px;
  left: 385px;
  z-index: 5;
  border-radius: 0 90px 90px 0;
}

#mirror-top {
  background: #e74c3c;
  width: 13px;
  height: 23px;
  top: 122px;
  left: 385px;
  z-index: 5;
  border-radius: 0 90px 90px 0;
}

#driver-bg {
  width: 68px;
  height: 29px;
  background: #222;
  top: 155px;
  left: 331px;
  z-index: 5;
  border-radius: 5px;
}

#driver-wheel {
  width: 5px;
  height: 25px;
  background: #95a5a6;
  top: 157px;
  left: 391px;
  z-index: 5;
  border-radius: 5px;
}

#bottom-body-spine {
  background: #c0392b;
  width: 80px;
  height: 5px;
  top: 197px;
  left: 300px;
  z-index: 4;
  opacity: 0.2;
}

#top-body-spine {
  background: #c0392b;
  width: 80px;
  height: 5px;
  top: 135px;
  left: 300px;
  z-index: 4;
  opacity: 0.2;
}

#bottom-body-spine-2 {
  background: #c0392b;
  width: 115px;
  height: 5px;
  top: 187px;
  left: 186px;
  z-index: 4;
  transform: rotate(10deg);
  opacity: 0.2;
}

#top-body-spine-2 {
  background: #c0392b;
  width: 115px;
  height: 5px;
  top: 145px;
  left: 186px;
  z-index: 4;
  transform: rotate(-10deg);
  opacity: 0.2;
}

#end-body-spine {
  width: 80px;
  height: 30px;
  top: 153px;
  left: 180px;
  z-index: 4;
  border-radius: 50px 0px 0px 50px;
  border-left: 6px solid #c0392b;
  opacity: 0.2;
}

#driver-helmet {
  background: #3498db;
  width: 27px;
  height: 25px;
  top: 157px;
  left: 332px;
  z-index: 6;
  border-radius: 20px 50px 50px 20px;
}

/* Responsive design */
@media (max-width: 768px) {
  #Container {
    transform: scale(0.3);
  }
  
  .loading-text h2 {
    font-size: 1.1rem;
  }
  
  .loading-text {
    max-width: 200px;
  }

  .loading-text p {
    font-size: 0.8rem;
  }

  .preloader-content {
    gap: 1.2rem;
  }
}

@media (max-width: 480px) {
  #Container {
    transform: scale(0.25);
  }
  
  .loading-text h2 {
    font-size: 1rem;
  }

  .loading-text p {
    font-size: 0.75rem;
  }

  .preloader-content {
    gap: 1rem;
  }
}
</style>