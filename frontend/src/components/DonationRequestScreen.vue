<template>
  <div>
    <h2>Request a Meal Donation</h2>
    <div v-if="isAvailableMealTypesPending">
      <p>Loading available meals...</p>
    </div>
    <div v-else-if="isAvailableMealTypesError">
      <p>Error loading meals. Please try again later.</p>
    </div>
    <div v-else class="flex">
      <p>Select meals you'd like to receive. Matching donations are assigned in request order (earliest requests fulfilled first).</p>
      <div class="flex-left full-width">
        <InputText
            id="requester-name"
            v-model="name"
            :invalid="userNameInputError !== ''"
            class="full-width"
            placeholder="Name"
        />
        <small v-if="userNameInputError !== ''" class="error-text">{{ userNameInputError }}</small>
      </div>
      <div class="flex-left full-width">
        <Listbox
            v-model="selectedMeals"
            :options="mealPreferencesComputed"
            optionLabel="description"
            :invalid="mealPreferencesError !== ''"
            multiple
            :highlightOnSelect="false"
            checkmark
        />
        <small v-if="mealPreferencesError !== ''" class="error-text">{{ mealPreferencesError }}</small>
      </div>
      <div class="button-row">
        <Button
            class="m-right"
            @click="goBack"
            outlined
        >
          Cancel
        </Button>
        <Button
            @click="submitDonationRequest"
            :disabled="isSubmittingRequest"
        >
          Submit Request
        </Button>
      </div>
    </div>
    <Dialog :visible="requestSuccessDialogVisible" header="Request Submitted!" modal>
      <p>Your meal donation request has been submitted. Check back on the "Receive Meal" page to see if a meal has been received.</p>
      <template #footer>
        <Button label="Okay" @click="handleRequestOkayButton" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useQuery, useMutation } from '@tanstack/vue-query';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import Listbox from "primevue/listbox";
import { useToast } from 'primevue/usetoast';
import api from '../axios/axios';
import { getNameFromCookie, setNameCookie } from '../utils/utils';
import type { ApiResult, Meal, MealPreference } from '../models/models';

const router = useRouter();
const toast = useToast();

const name = ref(getNameFromCookie() || '');
const requestSuccessDialogVisible = ref(false);
const userNameInputError = ref('');
const mealPreferences = ref<MealPreference[]>([]);
const mealPreferencesError = ref('');
const isSubmittingRequest = ref(false);
const selectedMeals = ref<MealPreference[]>([]);

const {isPending: isAvailableMealTypesPending, data: availableMealTypesData, isError: isAvailableMealTypesError} = useQuery({
  queryKey: ['availableMealTypes'],
  queryFn: async (): Promise<Meal[]> => {
    try {
      const response = await api.get('/Api/Meal/Today');
      const result: ApiResult<Meal[]> = response.data;
      return result.data || [];
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal types', life: 3000 });
      throw error;
    }
  },
  refetchOnWindowFocus: false,
});

const mealPreferencesComputed = computed(() => {
  if (availableMealTypesData.value) {
    mealPreferences.value = availableMealTypesData.value.map(meal => ({
      id: meal.id,
      description: meal.description,
      selected: false
    }));
  }
  return mealPreferences.value;
});

const donationRequestMutation = useMutation({
  mutationFn: async ({ requesterName, mealIds }: { requesterName: string, mealIds: number[] }) => {
    console.log('Submitting donation request:', { requesterName, mealIds });
    return await api.post('/Api/DonationRequest', {
      requesterName: requesterName,
      mealIds: mealIds
    });
  },
  onSuccess: () => {
    isSubmittingRequest.value = false;
    requestSuccessDialogVisible.value = true;
  },
  onError: (error) => {
    isSubmittingRequest.value = false;
    console.error('Donation request error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to submit donation request', life: 3000 });
  }
});

const validateDonationRequest = (userName: string): boolean => {
  let valid = true;

  if (!userName || userName.trim() === '') {
    userNameInputError.value = 'Please enter a name';
    valid = false;
  } else if (!/^(\w+\s?){1,5}$/.test(userName)) {
    userNameInputError.value = 'Please enter a valid name';
    valid = false;
  } else {
    userNameInputError.value = '';
  }

  const selectedMealIds = selectedMeals.value
    .map(meal => meal.id);

  if (selectedMealIds.length === 0) {
    mealPreferencesError.value = 'Please select at least one meal type';
    valid = false;
  } else {
    mealPreferencesError.value = '';
  }

  return valid;
};

const submitDonationRequest = () => {
  if (!validateDonationRequest(name.value)) {
    return;
  }

  setNameCookie(name.value);
  isSubmittingRequest.value = true;

  const selectedMealIds = selectedMeals.value
      .map(meal => meal.id)

  donationRequestMutation.mutate({
    requesterName: name.value,
    mealIds: selectedMealIds
  });
};

const goBack = () => {
  router.push('/receive-meal');
};

const handleRequestOkayButton = () => {
  requestSuccessDialogVisible.value = false;
  router.push('/');
};
</script>

<style scoped>
  .full-width {
    width: 100%;
  }

  .flex {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    justify-content: center;
    align-items: center;
  }

  .flex-left {
    display: flex;
    flex-direction: column;
    justify-content: left;
  }

  .error-text {
    text-align: left;
  }

  .meal-preference {
    display: flex;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .button-row {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    width: 100%;
  }

  .m-right {
    margin-right: 0.5rem;
  }

  .ml-2 {
    margin-left: 0.5rem;
  }
</style>