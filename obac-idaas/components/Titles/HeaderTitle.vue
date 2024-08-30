<template>
  <a :href="hrefs" class="group transition">
    <span
      class="font-orbitron text-neon-green transition group-hover:text-neon-green"
    >
      {{ sectionNumber }}.
    </span>
    <span class="text-gray-100 transition group-hover:text-neon-green">
      {{ sectionName }}</span
    >
  </a>
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
  hrefs: {
    type: String,
    required: false,
  },
});
</script>
