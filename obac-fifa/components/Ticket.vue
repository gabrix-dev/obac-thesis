<template>
  <div class="rounded-lg bg-white p-5 shadow-xl">
    <h2 class="mb-3 text-3xl font-bold">{{ ticket.ticketType }}</h2>
    <p class="mb-3 overflow-hidden text-sm">
      {{ ticket.ticketDescription }}
    </p>
    <h3 class="mb-3 text-2xl font-bold">{{ ticket.ticketPrice }}</h3>
    <button
      role="link"
      @click="openInNewTab(ticket.ticketBuyUrl)"
      class="ml-auto rounded-full bg-red-800 px-4 py-2 text-white hover:bg-red-400"
    >
      BUY TICKET
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  ticket: {
    type: Object,
    required: true,
  },
});
</script>

<script>
export default {
  methods: {
    openInNewTab(url) {
      window.open(url, "_blank", "noreferrer");
    },
  },
};
</script>
