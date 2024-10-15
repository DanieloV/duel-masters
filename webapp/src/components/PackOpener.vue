<template>
  <div class="flex flex-col align-center h-screen bg-gray-500 overflow-hidden">
    <div class="flex flex-row justify-center mt-3 flex-initial">
      <div @click="openPack()" class="btn save">Open Pack</div>
    </div>
    <div v-if="!opened" class="flex flex-col justify-end grow min-h-0 min-w-0 h-full items-center">
      <div class="h-3/6 w-max min-h-0 min-w-0 relative mb-8">
        <PackCard 
          v-for="(card, i) in cards"
          ref="cardsRef"
          :card="card"
          :style="{
              'top': i * 1 + 'px',
              'right': i * 1 + 'px',
          }"
          class="absolute h-full min-h-0"
        />
        <img class="absolute h-full min-h-0 z-30" src="/assets/images/packwrapper.webp" />
      </div>
    </div>
    <div 
      class="flex flex-row justify-around grow min-h-0 min-w-0 py-1"
      :class="{'hidden' : !opened}"  
    >
      <template v-for="(card, i) in cards">
        <PackCard v-if="i < 5" :card="card" ref="cardsRow1" :flipped="flipped[i]"/>
      </template>
    </div>
    <div class="flex flex-row justify-around grow min-h-0 min-w-0 py-1">
      <template v-for="(card, i) in cards">
        <PackCard v-if="i >= 5" :card="card" ref="cardsRow2" :flipped="flipped[i]"/>
      </template>
    </div>
  </div>
</template>

<script>
import PackCard from './PackCard.vue';

export default {

props: ['packs'],
components: {
  PackCard,
},
data() {
  return {
    opened: false,
    flipped: {
      0: false,
      1: false,
      2: false,
      3: false,
      4: false,
      5: false,
      6: false,
      7: false,
      8: false,
      9: false,
    }
  }
},

computed: {
  cards() {
    if (!this.packs) {
      return [] 
    }
    return this.packs[0]
  } 
},
methods: {
  openPack() {
    let rMap = {};
    this.$refs.cardsRef.forEach((cardRef, i) => {
      rMap[i] = {}
      rMap[i]["before"] = cardRef.$el.getBoundingClientRect()
    })
    this.opened = true
    this.$nextTick(() => {
      this.$refs.cardsRow1.forEach((cardRef, i) => {
        rMap[i]["after"] = cardRef.$el.getBoundingClientRect()
        rMap[i]["el"] = cardRef.$el
      })
      this.$refs.cardsRow2.forEach((cardRef, i) => {
        rMap[i + 5]["after"] = cardRef.$el.getBoundingClientRect()
        rMap[i + 5]["el"] = cardRef.$el
      })

      console.log(rMap)
      let biggestDist = 0

      Object.keys(rMap).forEach(key => {
        const deltaX = rMap[key]["before"].left - rMap[key]["after"].left;
        const deltaY = rMap[key]["before"].top - rMap[key]["after"].top;

        rMap[key]['deltaX'] = deltaX
        rMap[key]['deltaY'] = deltaY
        let deltaDis = Math.sqrt(deltaX**2 + deltaY**2)
        rMap[key]['dist'] = deltaDis
        if (deltaDis > biggestDist) {
          biggestDist = deltaDis
        }
      })
      
      Object.keys(rMap).forEach(key => {
        let time = rMap[key]['dist'] * 1500 / biggestDist

        let flipTime = 200 + key * 50
        
        setTimeout(() => {
          this.flipped[key]=true
        }, flipTime);
        

        rMap[key]["el"].animate([{
          transformOrigin: 'top left',
          transform: `
            translate(${rMap[key]["deltaX"]}px, ${rMap[key]["deltaY"]}px)
          `
        }, {
          transformOrigin: 'top left',
          transform: 'none'
        }], {
          duration: 1000,
          easing: 'ease-in-out',
        });
      })
    })
  }
}

}
</script>

<style scoped lang="scss">
.btn {
  display: inline-block;
  background: #5865F2;
  color: #e3e3e5;
  line-height: 20px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: 0.1s;
  text-align: center !important;
  user-select: none;
}

.btn:hover {
  cursor: pointer;
  background: #515de2;
}

.btn:active {
  background: #4c58d3 !important;
}

.save {
  background: #3ca374 !important;
  display: inline-block;
}

.save:hover {
  background: #3ca374 !important;
}

.card-proportions {
  aspect-ratio: 264/367;
}

</style>