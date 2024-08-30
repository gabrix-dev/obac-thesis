<template>
  <div class="mt-2">
    <button
      @click="$emit('accordionToggled')"
      class="flex w-full items-center justify-between rounded-lg bg-purple-300 py-2 px-4 hover:bg-purple-200"
    >
      <span class="text-lg font-medium">{{ props.title }}</span>
      <span v-if="props.show" class="text-lg">&#x25BC;</span>
      <span v-else class="text-lg">&#x25B6;</span>
    </button>
    <div v-if="props.show" class="w-full rounded-lg bg-purple-50 p-4 shadow-lg">
      <slot />
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  show: {
    type: Boolean,
    required: true,
  },
});
const emit = defineEmits(["accordionToggled"]);

//await getWhitelistedDbs();
</script>

<style>
button:focus {
  outline: none;
}
</style>
