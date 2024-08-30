<template>
  <div class="mx-auto max-w-md rounded-lg bg-white shadow-lg">
    <div
      class="flex items-center rounded-t-lg bg-green-500 py-4 px-6 text-white"
    >
      <i class="fas fa-check-circle mr-2"></i>
      <h2 class="text-lg font-semibold">Booking Confirmed</h2>
    </div>
    <div class="p-6">
      <div class="mb-4 flex items-center justify-between">
        <div class="flex items-center">
          <i class="fas fa-map-marker-alt mr-2 text-gray-600"></i>
          <div>
            <div class="text-sm text-gray-600">Zone:</div>
            <div class="text-2xl font-bold">{{ bookingInfoObject.zone }}</div>
          </div>
        </div>
        <div class="flex items-center">
          <i class="fas fa-chair mr-2 text-gray-600"></i>
          <div>
            <div class="text-sm text-gray-600">Seat:</div>
            <div class="text-2xl font-bold">{{ bookingInfoObject.seat }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  bookingInfoObject: {
    type: Object,
    required: true,
  },
});
console.log("ENTERING PROPS");
</script>
