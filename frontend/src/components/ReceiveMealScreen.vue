<template>
  <div>
    <h2>Receive a Meal</h2>
    <div v-if="isChosenMealsPending || isRequestSubmittedPending || isMealsPending">
      <p>Loading data...</p>
    </div>
    <div v-else-if="!isChosenMealsError && chosenMealsData">
      <p>You have selected "{{chosenMealsData.description}}" from {{chosenMealsData.donorName}}</p>
    </div>
    <div v-else-if="!isRequestSubmittedError && requestSubmittedData && requestSubmittedData.length > 0">
      <p>No meals have been donated yet that match your preferences. Try coming back later!</p>
    </div>
    <div v-else-if="isMealsError">
      <p>Error loading meals. Please try again later.</p>
    </div>
    <div v-else-if="mealsData && mealsData.length === 0">
      <p>There are no meals available at the moment.</p>
      <Button
          class="full-width"
          @click="navigateToDonationRequest"
      >
        Request a Meal Donation
      </Button>
    </div>
    <div v-else-if="mealsData && mealsData.length > 0" class="flex">
      <div class="flex-left full-width">
        <InputText
            id="name"
            v-model="name"
            :invalid="userNameInputError !== ''"
            class="full-width"
            placeholder="Name"
            disabled
        />
        <small v-if="userNameInputError !== ''" class="error-text">{{ userNameInputError }}</small>
      </div>
      <div class="flex-left full-width">
        <Listbox
            v-model="selectedDonation"
            :invalid="mealInputError !== ''"
            :options="mealsData"
            optionLabel="description"
            class="full-width"
        />
        <small v-if="mealInputError !== ''" class="error-text">{{ mealInputError }}</small>
      </div>
      <Button
          class="full-width"
          @click="selectMeal"
      >
        Select Option
      </Button>
    </div>
    <Dialog :visible="dialogVisible" header="Meal Claimed!" modal>
      <p v-if="selectedDonation">You have claimed "{{ selectedDonation.description }}" from {{ selectedDonation.donorName }}</p>
      <template #footer>
        <Button label="Okay" @click="handleOkayButton" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue';
import {useRouter} from 'vue-router';
import {useMutation, useQuery, useQueryClient} from '@tanstack/vue-query';
import Listbox from 'primevue/listbox';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import {useToast} from 'primevue/usetoast';
import api from '../axios/axios';
import {setNameCookie} from '../utils/utils';
import type {ApiResult, Donation} from '../models/models';
import { userStore } from '../store/user';

const router = useRouter();
const toast = useToast();
const queryClient = useQueryClient();

const name = ref(userStore.user?.name || '');
const selectedDonation = ref<Donation>({} as Donation);
const dialogVisible = ref(false);
const userNameInputError = ref('');
const mealInputError = ref('');

const {isPending: isChosenMealsPending, data: chosenMealsData, isError: isChosenMealsError} = useQuery({
  queryKey: ['chosenMeal'],
  queryFn: async (): Promise<Donation | null> => {
    try {
      const response = await api.get(`/Api/Donation/Claim?name=${name.value}&timestamp=${new Date().getTime()}`);
      const result: ApiResult<Donation> = response.data;
      return result.data;
    } catch (error: any) {
      if (error.response.status == 404) {
        return null;
      }
      throw error;
    }
  },
  refetchOnWindowFocus: true,
});

const {isPending: isRequestSubmittedPending, data: requestSubmittedData, isError: isRequestSubmittedError} = useQuery({
  queryKey: ['requestSubmitted'],
  queryFn: async (): Promise<Donation[] | null> => {
    try {
      const encodedName = encodeURIComponent(name.value);
      const response = await api.get(`/Api/DonationRequest/User?name=${encodedName}&date=${new Date().toISOString().split('T')[0]}`);
      const result: ApiResult<Donation[]> = response.data;
      return result.data || [];
    } catch (error: any) {
      if (error.response.status == 404) {
        return null;
      }
      throw error;
    }
  },
  refetchOnWindowFocus: true,
});

const {isPending: isMealsPending, data: mealsData, isError: isMealsError} = useQuery({
  queryKey: ['availableMeals'],
  queryFn: async (): Promise<Donation[]> => {
    try {
      const response = await api.get(`/Api/Donation?timestamp=${new Date().getTime()}`);
      const result: ApiResult<Donation[]> = response.data;
      return result.data || [];
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
      throw error;
    }
  },
  refetchOnWindowFocus: true,
  staleTime: 60000
});

const claimMutation = useMutation({
  mutationFn: async ({ donationId, name }: { donationId: number, name: string }) => {
    return await api.post('/Api/Donation/Claim', {
      donationId,
      name
    });
  },
  onSuccess: () => {
    dialogVisible.value = true;
    queryClient.invalidateQueries({ queryKey: ['availableMeals'] });
  },
  onError: () => {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000 });
    queryClient.invalidateQueries({ queryKey: ['availableMeals'] });
  }
});

const validateDonationClaim = (userName: string, donation: Donation): boolean => {
  let valid = true;

  const id = donation?.id;
  const description = donation?.description;

  if (!userName || userName.trim() === '') {
    userNameInputError.value = 'Please enter a name';
    valid = false;
  } else {
    userNameInputError.value = '';
  }

  if (id == null || id <= 0 || description == null || description.trim() === '') {
    mealInputError.value = 'Please select a meal';
    valid = false;
  } else {
    mealInputError.value = '';
  }

  return valid;
};

const selectMeal = () => {
  if (!validateDonationClaim(name.value, selectedDonation.value)) {
    return;
  }

  setNameCookie(name.value);

  claimMutation.mutate({
    donationId: selectedDonation.value.id,
    name: name.value
  });
};

const handleOkayButton = () => {
  router.push('/');
};

const navigateToDonationRequest = () => {
  router.push('/donation-request');
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
</style>
