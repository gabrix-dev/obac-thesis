<template>
  <div class="flex">
    <div class="mt-4 flex w-full rounded-md border px-5 py-3">
      <div class="">
        <h2 class="mb-2 font-bold">Scope name</h2>
        <div class="flex">
          <InputsText
            hint="Scope name"
            :defaultValue="store.scopes[scopeIndex].name"
            @input="onScopeNameChange"
            class="h-9"
            variable-width
          />
          <ButtonsAdd @click="onAddButtonClick" class="mx-3 h-9" />
        </div>
        <div class="justify-end">
          <button
            type="button"
            id="applyPolicyButton"
            class="mt-4 rounded-lg border border-entrust-purple bg-gray-100 py-2 px-4 text-sm text-entrust-purple hover:bg-white"
            @click="onDeleteButtonClick"
          >
            Delete scope
          </button>
        </div>
      </div>

      <div class="ml-7 w-full">
        <div
          class="flex items-center"
          v-for="(claim, index) in store.scopes[scopeIndex].claims"
          :key="index"
        >
          <div class="w-1/3">
            <h2 v-if="index == 0" class="mb-2 font-bold">Association type</h2>
            <InputsSelect
              @input="(value) => onTypeChange(value, index)"
              :options="associationTypes"
              :defaultValue="associationTypesAlias[claim.type]"
              class="mr-3"
            />
          </div>
          <div class="w-1/3">
            <h2 v-if="index == 0" class="mb-2 font-bold">Association value</h2>
            <InputsSelect
              @input="(value) => onValueChange(value, index)"
              :options="props.associationMapping[claim.type]"
              :defaultValue="claim.value"
            />
          </div>
          <div class="ml-2 w-1/3">
            <h2 v-if="index == 0" class="mx-2 font-bold">OpenID Claim</h2>
            <InputsText
              variable-width
              class="mx-2 mt-1 h-8"
              :defaultValue="claim.oidc"
              @input="(value) => onOidcChange(value, index)"
            />
          </div>
          <div>
            <IconsTrash
              class="ml-4 rounded-xl text-entrust-purple shadow-xl hover:cursor-pointer"
              :class="[index == 0 ? 'mt-9' : 'mt-0']"
              @click="onTrashButtonClick(index)"
            />
          </div>
        </div>
      </div>
    </div>
    <!-- <IconsTrash
      class="my-auto ml-2 rounded-xl text-entrust-purple shadow-xl hover:cursor-pointer"
      @click="onTrashButtonClick(index)"
    /> -->
  </div>
</template>

<script setup>
import { store } from "~/store.js";
const props = defineProps({
  associationMapping: {
    type: Object,
    required: true,
  },
  scopeIndex: {
    type: Number,
    required: true,
  },
});

var timerId;
var scopeNameInput;
const associationTypes = ["", "Asset Attribute Value", "Asset Related Value"];
const associationTypesAlias = {
  asset_attribute: "Asset Attribute Value",
  asset_related: "Asset Related Value",
};
var name = store.getName(props.scopeIndex);

console.log(
  "Association mapping prop: " + JSON.stringify(props.associationMapping),
);

function onTypeChange(value, claimIndex) {
  // console.log(
  //   "Type change: old - new: " +
  //     store.getScope(props.scopeIndex).claims[claimIndex].type +
  //     "-" +
  //     value,
  // );
  let storeValue = "asset_related";
  if (value == "Asset Attribute Value") {
    storeValue = "asset_attribute";
  }
  store.setClaimType(storeValue, claimIndex, props.scopeIndex);
  console.log(
    "Current index claims: " +
      JSON.stringify(store.scopes[props.scopeIndex].claims),
  );
}

function onValueChange(value, claimIndex) {
  store.setClaimValue(value, claimIndex, props.scopeIndex);
}

function onAddButtonClick() {
  store.addDefaultClaim(props.scopeIndex);
}

function onScopeNameChange(value) {
  clearTimeout(timerId);

  timerId = setTimeout(() => {
    console.log("Name update: " + value);
    store.setScopeName(props.scopeIndex, value);
  }, 500);
}

function onTrashButtonClick(index) {
  store.deleteClaim(props.scopeIndex, index);
}

function onDeleteButtonClick() {
  store.deleteScope(props.scopeIndex);
}

function onOidcChange(value, index) {
  clearTimeout(timerId);

  timerId = setTimeout(() => {
    console.log("OIDC update: " + value);
    store.setClaimOidc(value, index, props.scopeIndex);
  }, 500);
}
</script>
