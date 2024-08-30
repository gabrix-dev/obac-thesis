<template>
  <div class="mt-4 w-full px-8">
    <PolicyInput
      hint="Name"
      :defaultValue="props.initialValue.name"
      isMandatory="true"
      @policyUpdate="(value) => (ruleName = value)"
    >
      Name:
    </PolicyInput>
    <InputsRadio
      @input="(value) => (currentBC = value)"
      :initialValue="
        props.initialValue.ownershipConstraints[0].smartContract.blockchain
      "
      :id="id"
    />
    <!-- <div class="mt-8 flex">
      <h3 class="pr-5 text-lg">Action:</h3>
      <InputsSelect
        :options="['Permit', 'Deny']"
        @input="(selectValue) => (actionSelect = selectValue)"
      />
    </div> -->

    <div class="mt-8 flex">
      <h3 class="mr-4 pr-0 pb-5 text-lg">Whitelist:</h3>
      <div class="w-full">
        <div class="mb-0 flex">
          <InputsText
            hint="Add address"
            class="mr-5 h-8"
            @input="(value) => (addAddressInput = value)"
          ></InputsText>
          <ButtonsAdd @click="onAddAddressClick" class="mb-2" />
          <!-- <ButtonsMinus @click="onDeleteAddressClick" class="mb-2 ml-10" />
          <InputsText
            hint="Delete address"
            class="ml-5 h-8"
            @input="(value) => (deleteAddressInput = value)"
          ></InputsText> -->
        </div>
        <div class="mb-5 flex"></div>
      </div>
    </div>
    <ul class="rounded-lg border p-2" v-if="showWhitelistedBorder">
      <li
        v-for="(address, i) in whitelistedList"
        :class="[i % 2 == 0 ? 'bg-gray-100' : ' bg-white']"
        class="flex items-center justify-between"
      >
        {{ address }}
        <span class="">
          <IconsBucket @click="onDeleteWhitelistedClick(i)" />
        </span>
      </li>
    </ul>

    <PolicyInput
      hint="Smart contract address"
      :defaultValue="scAddress"
      @policyUpdate="(value) => (scAddress = value)"
      >Smart contract address:</PolicyInput
    >

    <div>
      <div class="mt-6 flex items-center space-x-4">
        <p class="text-lg">Asset constraints:</p>
        <ButtonsAdd @click="onAddButtonClick" />
        <ButtonsMinus @click="onMinusButtonClick" />
      </div>
      <div class="mt-7 h-44 overflow-y-auto rounded-md border-2">
        <InputsConstraint
          v-for="constrain in metadataConstraints"
          :value="constrain"
          @constrainUpdate="
            (constrainType, newValue, id) =>
              (metadataConstraints[id][constrainType] = newValue)
          "
        />
      </div>
    </div>

    <div v-if="showSuccessMessage" class="mt-4">
      <BadgesSuccess
        @closeBadgeButtonClicked="() => (showSuccessMessage = false)"
        >The rule has been updated</BadgesSuccess
      >
    </div>
    <BadgesError
      class="mt-4"
      v-if="showErrorMessage"
      @closeBadgeButtonClicked="() => (showErrorMessage = false)"
      >{{ errorMessage }}</BadgesError
    >
    <div class="mx-auto mt-12 flex w-1/3 justify-center">
      <button
        type="button"
        id="applyPolicyButton"
        class="mb-8 w-full rounded-lg bg-entrust-purple py-2 px-4 text-white hover:bg-purple-500"
        @click="onSaveRuleButtonClick"
      >
        Save rule
      </button>
      <button
        type="button"
        id="applyPolicyButton"
        class="mb-8 ml-3 w-full rounded-lg bg-red-600 py-2 px-4 text-white hover:bg-red-500"
        @click="onDeleteButtonClick"
      >
        Delete rule
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
const emit = defineEmits(["saveRule"]);
const props = defineProps({
  initialValue: {
    type: Object,
    required: true,
  },
  id: {
    type: Number,
    required: true,
  },
});

console.log("NEW RULE WITH ID: " + props.id);

const config = useRuntimeConfig();
const idpBackUrl = config.public.IDP_SERVER_URL;

const showSuccessMessage = ref(false);
const showErrorMessage = ref(false);
const isWaitingResponse = ref(false);
const whitelistedList = ref([]);
const metadataConstraints = ref([]);
const showWhitelistedBorder = ref(false);
var currentBC =
  props.initialValue.ownershipConstraints[0].smartContract.blockchain;
var actionSelect = "Permit";
var ruleName = props.initialValue.name;
var scAddress = "";
var addAddressInput = ref("");
var deleteAddressInput = ref("");
var errorMessage;

setInitialInputValues();

function setInitialInputValues() {
  console.log("Initial values: " + JSON.stringify(props.initialValue));
  if (props.initialValue.accessAddresses != undefined) {
    for (let address of props.initialValue.accessAddresses) {
      whitelistedList.value.push(address);
      if (whitelistedList.value.length == 1) {
        showWhitelistedBorder.value = true;
      }
    }
  }
  addAddressInput.value = "";
  if (props.initialValue.ownershipConstraints[0] != undefined) {
    scAddress =
      props.initialValue.ownershipConstraints[0].smartContract.address;
    for (let metadataConstraint of props.initialValue.ownershipConstraints[0]
      .metadataConstraints) {
      metadataConstraints.value.push({
        key: metadataConstraint.key,
        value: metadataConstraint.value.join(","),
        operator: metadataConstraint.operator,
        id: metadataConstraints.value.length,
      });
    }
  }
}

function arrayToString() {
  return arr.join(",");
}

function onAddAddressClick() {
  whitelistedList.value.push(addAddressInput.value);
  if (whitelistedList.value.length == 1) {
    showWhitelistedBorder.value = true;
  }
  // addAddressInput.value = ""; no va aixi --> Todo: afegir watcher per poder borrar inputs
}

function onDeleteWhitelistedClick(i) {
  console.log("click!");
  if (whitelistedList.value.length == 1) {
    showWhitelistedBorder.value = false;
  }
  whitelistedList.value.splice(i, 1);
}

function onAddButtonClick() {
  console.log("Constrain list: " + JSON.stringify(metadataConstraints.value));
  metadataConstraints.value.push({
    key: "",
    value: "",
    operator: "Equals",
    id: metadataConstraints.value.length,
  });
}

function transformConstrainValue(value) {
  let arr = value.split(",").map((item) => item.trim());
  return arr;
}

function onMinusButtonClick() {
  if (metadataConstraints.value.length > 1) {
    metadataConstraints.value.pop();
  }
}

function buildRuleObject() {
  const rule = {
    name: ruleName,
    action: actionSelect,
  };

  if (whitelistedList.value.length > 0) {
    rule["accessAddresses"] = whitelistedList.value;
  }

  if (scAddress != "") {
    rule["ownershipConstraints"] = [];
    rule["ownershipConstraints"][0] = {};
    rule["ownershipConstraints"][0]["smartContract"] = {
      address: scAddress,
      blockchain: currentBC,
    };
    if (metadataConstraints.value[0].key != "") {
      rule["ownershipConstraints"][0]["metadataConstraints"] = [];
      for (let constrain of metadataConstraints.value) {
        let updatedConstrain = Object.assign({}, constrain);
        updatedConstrain.value = transformConstrainValue(
          updatedConstrain.value,
        );
        delete updatedConstrain.id;
        rule["ownershipConstraints"][0]["metadataConstraints"].push(
          updatedConstrain,
        );
      }
    }
  }
  console.log("RULE: " + JSON.stringify(rule));
  return rule;
}

function onDeleteButtonClick() {
  emit("deleteRule");
}

function onSaveRuleButtonClick() {
  showErrorMessage.value = false;
  showSuccessMessage.value = true;
  const rule = buildRuleObject();
  emit("saveRule", rule);
}
</script>
