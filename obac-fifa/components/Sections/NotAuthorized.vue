<template>
  <div class="">
    <div class="rounded-lg bg-white p-10">
      <h2 class="mb-5 text-center text-xl font-medium">
        You are not authorized to access this page
      </h2>
      <p class="mb-10 text-center text-gray-600">
        You need to authenticate with your ticket to access.
      </p>
      <div class="flex justify-center">
        <NuxtLink
          to="/"
          class="ml-5 rounded bg-red-800 py-3 px-6 text-white hover:bg-red-400"
          >Go back to menu</NuxtLink
        >
      </div>
    </div>
  </div>
</template>
