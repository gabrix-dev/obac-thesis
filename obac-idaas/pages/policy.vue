<template>
  <SectionsIDaaSAdminPortal currentPage="policy">
    <div class="mt-4 w-full px-8">
      <p class="mb-10 pt-4 text-4xl font-bold text-entrust-purple">
        New Ownership Policy
      </p>
      <div v-if="isLoadingPolicy" class="">
        <AnimationsLoading class="mx-auto w-40 py-36" />
      </div>
      <div v-else>
        <PolicyInput
          hint="Name"
          :defaultValue="policyNameInput"
          isMandatory="true"
          @policyUpdate="(value) => (policyNameInput = value)"
        >
          Name:
        </PolicyInput>
        <PolicyInput
          hint="Name"
          :defaultValue="policyDescriptionInput"
          isMandatory="true"
          @policyUpdate="(value) => (policyDescriptionInput = value)"
        >
          Description:
        </PolicyInput>
        <div class="mb-7 mt-10 flex items-center">
          <h2 class="mr-6 mt-1 text-lg">Access Control Rules</h2>
          <InputsText
            hint="Rule name"
            class="mr-5"
            @input="(value) => (addRuleInput = value)"
          ></InputsText>
          <ButtonsAdd @click="onAddRuleClick" />
        </div>

        <div v-for="(rule, index) in store.rules" :key="index" class="mt-5">
          <Accordion
            :show="showRules[index]"
            :title="rule.name"
            @accordionToggled="onAccordionToggled(index)"
          >
            <SectionsPolicy
              @saveRule="(rule) => onSaveRule(index, rule)"
              @deleteRule="onDeleteRule(index)"
              :initialValue="rule"
              :id="index"
            />
          </Accordion>
        </div>
        <div class="flex w-full justify-center">
          <button
            type="button"
            id="applyPolicyButton"
            class="mb-8 mt-8 w-full rounded-lg bg-entrust-purple py-2 px-4 text-white hover:bg-purple-500"
            @click="onApplyButtonClick"
          >
            Apply policy
          </button>
        </div>
        <div v-if="showSuccessMessage" class="">
          <BadgesSuccess
            @closeBadgeButtonClicked="() => (showSuccessMessage = false)"
            >The policy has been updated</BadgesSuccess
          >
        </div>
        <BadgesError
          class="mt-4"
          v-if="showErrorMessage"
          @closeBadgeButtonClicked="() => (showErrorMessage = false)"
          >{{ errorMessage }}</BadgesError
        >
        <AnimationsLoading v-if="isWaitingResponse" class="mx-auto pb-7" />
      </div>
    </div>
  </SectionsIDaaSAdminPortal>
</template>

<script setup>
import { store } from "~/store.js";

var addRuleInput;
var policyNameInput = "Policy name";
var policyDescriptionInput = "OBAC test";
var showRules = ref([]);
var errorMessage;
const showSuccessMessage = ref(false);
const showErrorMessage = ref(false);
const isWaitingResponse = ref(false);
const isLoadingPolicy = ref(true);
const showFormWarning = ref(false);
const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;

initRules();

function initRules() {
  //store.addDefaultRule("Default rule");
  fetch(idpBackUrl + "/policy/getPolicy", {
    method: "GET",
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Something went wrong retreiving the rules");
      }
    })
    .then((data) => {
      console.log("RECEIVED RULES: " + JSON.stringify(data));
      store.setRules(data.policies.accessControlPolicy.accessControlRules);
      policyNameInput = data.metadata.name;
      policyDescriptionInput = data.metadata.description;
      showRules.value = new Array(store.rules.length).fill(false);
    })
    .catch((error) => {
      console.log(error);
    })
    .finally(() => {
      isLoadingPolicy.value = false;
    });
}

function onAddRuleClick() {
  store.addDefaultRule(addRuleInput);
  showRules.value.push(false);
  console.log("RULES updated: " + JSON.stringify(store.rules));
}

function onSaveRule(index, rule) {
  store.updateRule(index, rule);
}

function onDeleteRule(index) {
  console.log("PRE-delete rules: " + JSON.stringify(store.rules));
  console.log("DELETING rule with index: " + index);
  store.deleteRule(index);
  console.log("POST-delete rules: " + JSON.stringify(store.rules));
  showRules.value.splice(index, 1);
}

function onAccordionToggled(index) {
  showRules.value[index] = !showRules.value[index];
}

function buildPolicyObject() {
  let rules = store.getRules();
  let policy = {
    metadata: {
      name: policyNameInput,
      description: policyDescriptionInput,
    },
    policies: {
      accessControlPolicy: {
        accessControlRules: rules,
      },
    },
  };
  return policy;
}

function onApplyButtonClick() {
  showErrorMessage.value = false;
  showSuccessMessage.value = false;
  isWaitingResponse.value = true;
  let policyObj = buildPolicyObject();
  console.log("Policy JSON: " + JSON.stringify(policyObj));
  fetch(idpBackUrl + "/policy/edit", {
    method: "POST",
    mode: "cors",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(policyObj),
  })
    .then((response) => {
      if (response.ok) {
        console.log("Response: OK");
        isWaitingResponse.value = false;
        showSuccessMessage.value = true;
      } else {
        throw new Error("Invalid policy format");
      }
    })
    .catch((error) => {
      console.log(error);
      errorMessage = error.message;
      isWaitingResponse.value = false;
      showErrorMessage.value = true;
    });
}
</script>
