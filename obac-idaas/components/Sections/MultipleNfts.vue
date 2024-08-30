<template>
  <h2 class="mt-3 ml-2 mb-4 text-lg font-bold">
    Please select the asset you want to authenticate with:
  </h2>
  <div class="w-30 relative flex overflow-x-auto">
    <TicketImageId
      v-for="asset in assetUnits"
      :tokenId="asset.assetData.id"
      :collection="asset.collection"
      :blockchain="asset.blockchain"
      :imageUrl="asset.assetData.imageUrl"
      :isSelected="
        selectedId + selectedCollection + selectedBlockchain ==
        asset.assetData.id + asset.collection + asset.blockchain
      "
      class="p-2"
      @click="
        onAssetClick(asset.assetData.id, asset.collection, asset.blockchain)
      "
    />
  </div>
</template>

<script setup>
const emit = defineEmits(["assetClicked"]);
const selectedId = ref(-1);
const selectedCollection = ref("");
const selectedBlockchain = ref("");
const props = defineProps({
  userCollectionsData: {
    type: Array,
    required: true,
  },
});
var assetUnits = [];

onRender();

function onRender() {
  for (let userCollectionData of props.userCollectionsData) {
    const updatedList = userCollectionData.assetData.map((assetData) => {
      const collection = userCollectionData.name;
      const blockchain =
        userCollectionData.blockchain.charAt(0).toUpperCase() +
        userCollectionData.blockchain.slice(1);
      return { assetData, collection, blockchain };
    });
    assetUnits = assetUnits.concat(updatedList);
  }
}

function onAssetClick(id, collection, blockchain) {
  selectedId.value = id;
  selectedCollection.value = collection;
  selectedBlockchain.value = blockchain;
  emit("assetClicked", id, collection, blockchain);
}
</script>
