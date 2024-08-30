<template>
  <div
    class="transform rounded-lg px-6 py-4 shadow-xl transition-colors duration-200 hover:bg-gray-200"
  >
    <p class="text-lg font-medium text-gray-800 dark:text-gray-100">
      {{ ticket.ticketType }}
    </p>
    <h4 class="mt-2 text-4xl font-semibold text-gray-800 dark:text-gray-100">
      {{ ticket.ticketPrice }}
    </h4>
    <p class="mt-4 text-gray-500 dark:text-gray-300">
      {{ ticket.ticketDescription }}
    </p>

    <div class="mt-8 space-y-8">
      <ListItem v-for="bullet in ticket.ticketBulletPoints">
        {{ bullet }}</ListItem
      >
    </div>

    <button
      class="mt-10 w-full transform rounded-md bg-red-800 px-4 py-2 font-medium capitalize tracking-wide text-white focus:bg-red-400 focus:outline-none hover:bg-red-400"
      :class="[
        !walletConnected ? 'cursor-not-allowed opacity-50' : 'opacity-100',
      ]"
      @click="$emit('buyButtonClicked', ticket.ticketType)"
    >
      Buy
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  ticket: {
    type: Object,
    required: true,
  },
  walletConnected: {
    type: Boolean,
    required: true,
  },
});
</script>
