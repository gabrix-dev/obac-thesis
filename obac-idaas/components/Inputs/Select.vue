<template>
  <div>
    <select
      id="location"
      name="location"
      class="mt-1 block w-full rounded-md border-gray-300 py-2 pl-1 pr-2 text-base focus:border-purple-500 focus:outline-none focus:ring-purple-500 sm:text-sm"
      @input="$emit('input', $event.target.value)"
    >
      <option
        :selected="props.defaultValue == option"
        v-for="option in options"
      >
        {{ option }}
      </option>
    </select>
  </div>
</template>

<script setup>
defineEmits(["input"]);
const props = defineProps({
  options: {
    type: Array,
    required: true,
  },
  defaultValue: {
    type: String,
    required: false,
  },
});
</script>
