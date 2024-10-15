<template>
  <div class="card-container content-center max-h-full min-w-0"
    @click="flipped = true"
  >
    <div class="card"
      :class="{'flip': flipped}"
    >
      <img class="max-w-full max-h-full image-corners card-front" :src="`https://scans.shobu.io/${card}.jpg`" />
      <img class="max-w-full max-h-full image-corners card-back" :src="`https://scans.shobu.io/backside.jpg`" />
    </div>
  </div>
</template>

<script setup>
const flipped = defineModel('flipped')
const props = defineProps({
  card: String,
  opened: Boolean,
})

</script>

<style scoped lang="scss">
.card-container {
  aspect-ratio: 264/367;
  perspective: 1000px;
}

.flip {
  transform-style: preserve-3d;
  animation: rotateSuper 1s linear;
  animation-fill-mode: forwards ;    
}

@keyframes rotateSuper {
  50% {
    transform: rotate3d(1, 6, 2, -90deg);
  }
  100% {
    transform: rotate3d(0, 1, 0, -180deg);
  }
}

.image-corners {
  border-radius: 5% / 4%;
}

.card {
  position: relative;
  height: 100%;
  width: 100%;
}

.card-back, .card-front {
  position: absolute;
  display:  block;
  backface-visibility: hidden;
}

.card-front {
  transform: rotateY(180deg);
}

.card-back {
  
}


</style>