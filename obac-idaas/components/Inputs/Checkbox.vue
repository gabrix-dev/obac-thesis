<template>
  <div class="mr-4">
    <input
      type="checkbox"
      id="whitelisted"
      class="form-checkbox"
      @input="$emit('input', $event.target.value)"
    />
    <label for="whitelisted" class="ml-2 text-gray-700"><slot /></label>
  </div>
</template>

<script setup>
defineEmits(["input"]);
</script>

<script></script>
