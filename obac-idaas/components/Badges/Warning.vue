<template>
  <div
    id="alert-3"
    class="mb-4 flex items-center rounded-lg border border-yellow-400 bg-yellow-100 p-4"
    role="alert"
  >
    <IconsWarning class="h-4 w-4 text-yellow-700" />

    <div class="ml-2 text-sm font-bold text-yellow-600">Warning!</div>
    <div class="ml-1 text-sm text-yellow-600">
      <slot />
    </div>
    <button
      @click="$emit('closeBadgeButtonClicked')"
      type="button"
      class="-mx-1.5 -my-1.5 ml-auto inline-flex h-8 w-8 rounded-lg bg-yellow-100 p-1.5 text-yellow-600 focus:ring-2 focus:ring-yellow-400 hover:bg-yellow-400"
      data-dismiss-target="#alert-3"
      aria-label="Close"
    >
      <span class="sr-only">Close</span>
      <svg
        class="h-5 w-5"
        fill="currentColor"
        viewBox="0 0 20 20"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill-rule="evenodd"
          d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
          clip-rule="evenodd"
        ></path>
      </svg>
    </button>
  </div>
</template>

<script setup>
const emit = defineEmits(["closeBadgeButtonClicked"]);
</script>
