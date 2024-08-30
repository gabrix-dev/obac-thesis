<template>
    <select class="mt-6 border border-gray-400 rounded py-2 px-4">
      <option value="" disabled selected>Select a match</option>
      <option v-for="match in matches">{{ match }}</option>
    </select>
</template>

<script setup>
const props = defineProps({
  matches: {
    type: Array,
    required: true,
  }
});
</script>