<template>
  <div class="w-full">
    <img
      src="../../assets/img/stadium2.png"
      alt="Image 1"
      class="mx-auto w-1/5"
    />
    <img src="../../assets/img/zones.png" alt="Image 2" class="mx-auto w-1/3" />
    <h2 class="mt-6">Select a zone:</h2>
    <select
      class="mt-2 w-full rounded border border-gray-400 py-2 px-4"
      v-model="selectedZone"
      @change="$emit('zoneSelected', selectedZone)"
    >
      <option v-for="zone in zones">{{ zone }}</option>
    </select>
    <h2 class="mt-6">Select a seat:</h2>
    <select
      class="mt-2 w-full rounded border border-gray-400 py-2 px-4"
      v-model="selectedSeat"
      @change="$emit('seatSelected', selectedSeat)"
    >
      <option v-for="seat in seats">{{ seat }}</option>
    </select>
  </div>
</template>

<script setup>
const emit = defineEmits(["zoneSelected", "seatSelected"]);
const seats = [1, 2, 3, 4, 5, 6, 7, 8, 9];
const zones = ["VIP Zone", "Zone 1", "Zone 2", "Zone 3"];
</script>
