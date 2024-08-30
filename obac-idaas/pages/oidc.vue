<template>
  <SectionsIDaaSAdminPortal currentPage="oidc">
    <div class="mb-8 px-8">
      <div class="mt-4 w-full py-8">
        <p class="mb-10 pt-4 text-4xl font-bold text-entrust-purple">
          Manage OIDC ownership claims
        </p>
        <h2 class="mt-6 text-2xl text-gray-400">
          On this page, you can map ownership claims of the assets managed by
          the smart contracts in your policy to OIDC scopes
        </h2>
      </div>
      <div v-if="isLoadingMapping || isLoadingAssociations">
        <AnimationsLoading class="mx-auto w-40 py-52" />
      </div>
      <div v-else>
        <div class="flex items-center">
          <InputsText
            hint="Scope"
            @input="(value) => (scopeNameInput = value)"
            class="h-9"
          />
          <ButtonsAdd @click="onAddScopeButtonClick" class="mx-3 h-9" />
        </div>
        <SectionsScope
          v-for="(scope, index) in store.scopes"
          :key="index"
          :scopeIndex="index"
          :associationMapping="associationMapping"
        />
      </div>
    </div>
    <div class="mx-8 mb-8" v-if="!(isLoadingMapping && isLoadingAssociations)">
      <button
        type="button"
        id="saveButton"
        class="w-full rounded-lg bg-entrust-purple p-2 text-white hover:bg-entrust-dark-purple"
        @click="onSaveButtonClick"
      >
        Save
      </button>
      <div class="mt-6">
        <AnimationsLoading v-if="isWaitingResponse" class="mx-auto" />
        <BadgesSuccess
          v-if="showSuccessMessage"
          @closeBadgeButtonClicked="showSuccessMessage = false"
          >Scope mapping saved</BadgesSuccess
        >
      </div>
    </div>
  </SectionsIDaaSAdminPortal>
</template>

<script setup>
import { store } from "~/store.js";
const initialAssociationMapping = {
  "": [""],
  "Asset Attribute Value": [""],
  "Asset Related Value": [""],
};
var scopeNameInput;
const isWaitingResponse = ref(false);
const showSuccessMessage = ref(false);
const isLoadingMapping = ref(true);
const isLoadingAssociations = ref(true);
const associationMapping = ref(initialAssociationMapping);
const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;

getAssociations();
getMapping();

function getMapping() {
  fetch(idpBackUrl + "/api/oidc/getMapping", {
    method: "GET",
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Something went wrong retreiving the mapping");
      }
    })
    .then((data) => {
      console.log("RECEIVED mapping: " + JSON.stringify(data));
      store.setScopes(data);
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      isLoadingMapping.value = false;
    });
}

function getAssociations() {
  fetch(idpBackUrl + "/api/oidc/getAssociations", {
    method: "GET",
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Something went wrong retreiving the associations");
      }
    })
    .then((data) => {
      console.log("RECEIVED associations: " + JSON.stringify(data));
      associationMapping.value = data;
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      isLoadingAssociations.value = false;
    });
}

function onAddScopeButtonClick() {
  store.addScope(scopeNameInput);
}

function onSaveButtonClick() {
  isWaitingResponse.value = true;
  console.log("SCOPES: " + JSON.stringify(store.getScopes()));
  fetch(idpBackUrl + "/api/oidc/postMapping", {
    method: "POST",
    mode: "cors",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(store.getScopes()),
  })
    .then((response) => {
      if (response.ok) {
        console.log("Success");
        showSuccessMessage.value = true;
      } else {
        throw new Error("Something went wrong");
      }
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      isWaitingResponse.value = false;
    });
}
</script>
