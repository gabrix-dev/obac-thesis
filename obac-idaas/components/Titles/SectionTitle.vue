<template>
  <div class="mt-12 border-b border-neon-green text-3xl lg:text-5xl">
    <span class="font-orbitron text-neon-green"> {{ sectionNumber }}. </span>
    <span class="font-aldrich text-gray-100"> {{ sectionName }}</span>
  </div>
</template>

<script setup>
const props = defineProps({
  sectionName: {
    type: String,
    required: true,
  },
  sectionNumber: {
    type: String,
    required: true,
  },
});
</script>
