<template>
  <div>
    <h2>Receive a Meal</h2>
    <div v-if="!isChosenMealsError && !isChosenMealsPending && chosenMealsData">
      <p>You have selected "{{chosenMealsData.description}}" from {{chosenMealsData.donorName}}</p>
    </div>
    <div v-else-if="isMealsPending">
      <p>Loading available meals...</p>
    </div>
    <div v-else-if="isMealsError">
      <p>Error loading meals. Please try again later.</p>
    </div>
    <div v-else-if="mealsData && mealsData.length === 0">
      <p>There are no meals available at the moment.</p>
    </div>
    <div v-else class="flex">
      <div class="flex-left full-width">
        <InputText
            id="name"
            v-model="name"
            :invalid="userNameInputError !== ''"
            class="full-width"
            placeholder="Name"
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import Listbox from 'primevue/listbox';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import { useToast } from 'primevue/usetoast';
import api from '../axios/axios';
import { getNameFromCookie, setNameCookie } from '../utils/utils';
import type { ApiResult, Donation } from '../models/models';

const router = useRouter();
const toast = useToast();
const queryClient = useQueryClient();

const name = ref(getNameFromCookie() || '');
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
  refetchOnWindowFocus: false,
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
  } else if (!/^(\w+\s?){1,5}$/.test(userName)) {
    userNameInputError.value = 'Please enter a valid name';
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

// Lifecycle hooks
onMounted(() => {
  // Initial name is already set in ref initialization
});
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
