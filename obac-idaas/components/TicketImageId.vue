<template>
  <div
    class="relative m-2 rounded-md border-entrust-purple hover:cursor-pointer hover:border-2"
    :class="[isSelected ? 'border-2 ' : 'border-0']"
  >
    <img
      :src="imageUrl"
      class="mx-auto h-24 w-48 max-w-xs pb-1"
      :title="imageTitle"
    />
    <div
      class="absolute bottom-2 right-3"
      :title="blockchain"
      v-if="blockchain == 'Ethereum'"
    >
      <LogosEthereum class="w-4" />
    </div>
    <div
      class="absolute bottom-2 right-2"
      :title="blockchain"
      v-if="blockchain == 'Polygon'"
    >
      <LogosPolygon class="w-6" />
    </div>
    <div
      class="absolute bottom-2 right-2"
      :title="blockchain"
      v-if="blockchain == 'Bsc'"
    >
      <LogosBSC class="h-6 w-6" />
    </div>
    <h2 class="w-40 truncate font-bold">#{{ tokenId }} - {{ collection }}</h2>
  </div>
</template>

<script setup>
const props = defineProps({
  tokenId: {
    type: Number,
    required: true,
  },
  isSelected: {
    type: Boolean,
    required: false,
  },
  collection: {
    type: String,
    required: true,
  },
  imageUrl: {
    type: String,
    required: false,
  },
  blockchain: {
    type: String,
    required: false,
  },
});

const imageTitle = props.collection + " - " + props.blockchain;
</script>
