<template>
  <input
    type="text"
    :placeholder="props.hint"
    class="rounded-lg border border-gray-300 p-2"
    :class="[variableWidth ? 'w-full' : 'w-max']"
    @input="$emit('input', $event.target.value)"
    :value="defaultValue"
  />
</template>

<script setup="props">
defineEmits(["input"]);
const props = defineProps({
  hint: {
    type: String,
    required: false,
  },
  defaultValue: {
    type: String,
    required: false,
  },
  variableWidth: {
    type: Boolean,
    required: false,
  },
});
</script>
