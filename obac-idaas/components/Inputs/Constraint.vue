<template>
  <div class="ml-4 mt-6 flex items-center space-x-4">
    <InputsText
      hint="Key"
      :defaultValue="value.key"
      @input="(keyValue) => $emit('constrainUpdate', 'key', keyValue, value.id)"
    />
    <InputsSelect
      @input="
        (selectValue) =>
          $emit('constrainUpdate', 'operator', selectValue, value.id)
      "
      :options="['Equals', 'Not equals']"
    />
    <InputsText
      hint="Value"
      :defaultValue="value.value"
      @input="
        (valueValue) => $emit('constrainUpdate', 'value', valueValue, value.id)
      "
    />
  </div>
</template>

<script setup>
var defaultKey;
var defaultValue;
const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
});
defineEmits(["constrainUpdate"]);

// if (props.id == 0) {
//   defaultKey = "EventID";
//   defaultValue = "Qatar2022";
// }
</script>
