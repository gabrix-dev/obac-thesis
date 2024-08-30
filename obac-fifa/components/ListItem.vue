<template>
  <div class="flex items-center">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-5 w-5 text-red-500"
      viewBox="0 0 20 20"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
        clip-rule="evenodd"
      />
    </svg>
    <!-- <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-5 w-5 text-red-500"
      viewBox="0 0 20 20"
      fill="currentColor"
      v-if="isTick"
    >
      <path
        d="M10 8.586l3.293-3.293a1 1 0 011.414 0l.586.586a1 1 0 010 1.414L12.414 10l3.293 3.293a1 1 0 01-1.414 1.414L11 11.414l-3.293 3.293a1 1 0 01-1.414-1.414L8.586 10 5.293 6.707a1 1 0 011.414-1.414L10 8.586z"
      />
    </svg> -->
    <span class="mx-4 text-gray-500">
      <slot />
    </span>
  </div>
</template>
