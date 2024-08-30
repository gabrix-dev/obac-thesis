<template>
  <div class="mt-6 flex items-center space-x-4">
    <p class="text-lg">Blockchain:</p>
    <ul class="grid w-full gap-6 md:grid-cols-3">
      <li>
        <input
          type="radio"
          :checked="selected == 'ethereum'"
          :id="uid1"
          :name="groupName"
          value="ethereum"
          class="peer hidden"
          @input="onInput($event.target.value)"
          required
        />
        <label
          :for="uid1"
          class="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white text-gray-500 peer-checked:border-2 peer-checked:border-blue-600 peer-checked:text-blue-600 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:peer-checked:text-blue-500 dark:hover:bg-gray-700 dark:hover:text-gray-300"
        >
          <div class="block py-5 pl-3">
            <div class="w-full text-lg font-semibold">Ethereum</div>
          </div>
          <div class="flex h-12 w-12 items-center">
            <LogosEthereum class="p-3"></LogosEthereum>
          </div>
        </label>
      </li>
      <li>
        <input
          type="radio"
          :id="uid2"
          :name="groupName"
          value="bsc"
          class="peer hidden"
          :checked="selected == 'bsc'"
          @input="onInput($event.target.value)"
        />
        <label
          :for="uid2"
          class="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white text-gray-500 peer-checked:border-2 peer-checked:border-blue-600 peer-checked:text-blue-600 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:peer-checked:text-blue-500 dark:hover:bg-gray-700 dark:hover:text-gray-300"
        >
          <div class="block py-5 pl-3">
            <div class="w-full text-lg font-semibold">BSC</div>
          </div>
          <LogosBSC class="h-12 w-12 pr-2"></LogosBSC>
        </label>
      </li>
      <li>
        <input
          type="radio"
          :id="uid3"
          :name="groupName"
          value="polygon"
          :checked="selected == 'polygon'"
          class="peer hidden"
          @input="onInput($event.target.value)"
        />
        <label
          :for="uid3"
          class="inline-flex w-full cursor-pointer items-center justify-between rounded-lg border border-gray-200 bg-white text-gray-500 peer-checked:border-2 peer-checked:border-blue-600 peer-checked:text-blue-600 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:peer-checked:text-blue-500 dark:hover:bg-gray-700 dark:hover:text-gray-300"
        >
          <div class="block py-5 pl-3">
            <div class="w-full text-lg font-semibold">Polygon</div>
          </div>
          <div class="flex h-12 w-12 items-center px-2">
            <LogosPolygon />
          </div>
        </label>
      </li>
    </ul>
  </div>
</template>

<script setup>
//CANVIAR provider
const emit = defineEmits(["input"]);
const props = defineProps({
  initialValue: {
    type: String,
    required: false,
  },
  id: {
    type: Number,
    required: true,
  },
});
var selected, uid1, uid2, uid3, groupName;

onInit();

function onInit() {
  selected = props.initialValue;
  groupName = "provider" + props.id;
  uid1 = groupName + "1";
  uid2 = groupName + "2";
  uid3 = groupName + "3";
  console.log("GroupName: " + groupName);
  console.log("ID: " + props.id);
}

function onInput(value) {
  emit("input", value);
}
</script>
