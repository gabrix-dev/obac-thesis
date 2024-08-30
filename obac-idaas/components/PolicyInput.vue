<template>
  <div class="mt-6 flex items-center space-x-4">
    <p class="text-lg">
      <slot /> <span class="bold text-red-500" v-if="isMandatory">*</span>
    </p>
    <InputsText
      @input="(value) => $emit('policyUpdate', value)"
      :hint="hint"
      :defaultValue="defaultValue"
    />
  </div>
</template>

<script setup="props">
defineEmits(["policyUpdate"]);
const props = defineProps({
  hint: {
    type: String,
    required: false,
  },
  isMandatory: {
    type: Boolean,
    required: false,
  },
  defaultValue: {
    type: String,
    required: false,
  },
});
console.log("Prop hint: " + props.hint);
</script>
