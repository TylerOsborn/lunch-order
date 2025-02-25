<template>
  <div class="give-meal-screen">
    <h2>Give a Meal</h2>
    <form class="flex" @submit.prevent="submitMeal">
      <div class="flex-left full-width">
        <InputText
            class="full-width"
            :invalid="userNameInputErrorText !== ''"
            placeholder="Name"
            id="name"
            v-model="name"
        />
        <small v-if="userNameInputErrorText !== ''" id="name-help" class="error-text">
          {{ userNameInputErrorText }}
        </small>
      </div>

      <div class="flex-left full-width">
        <Listbox
            class="full-width"
            :class="{ 'invalid': mealInputErrorText !== '' }"
            v-model="selectedMealType"
            :options="meals"
            optionValue="id"
            optionLabel="description"
            placeholder="Select..."
            id="meal"
            required
        />
        <small v-if="mealInputErrorText !== ''" class="error-text">
          {{ mealInputErrorText }}
        </small>
      </div>

      <Button
          class="full-width"
          type="submit"
          :disabled="isPending"
      >
        Submit
      </Button>
    </form>
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useRouter } from 'vue-router';
import type { Meal, ApiResult } from '../models/models';
import { getNameFromCookie, setNameCookie } from '../utils/utils';
import Listbox from 'primevue/listbox';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import api from "../axios/axios.ts";
import { useToast } from 'primevue/usetoast';

const toast = useToast();


const router = useRouter();
const queryClient = useQueryClient();

const name = ref(getNameFromCookie());
const selectedMealType = ref(0);
const userNameInputErrorText = ref('');
const mealInputErrorText = ref('');

const { isPending, data: mealsResult } = useQuery({
  queryKey: ['meals', 'today'],
  queryFn: async (): Promise<ApiResult<Meal[]>> => {
    const response = await api.get('/Api/Meal/Today');
    return response.data;
  }
});

const meals = computed(() => mealsResult.value?.data || []);

const donationMutation = useMutation({
  mutationFn: async (donation: { donorName: string; mealId: number }) => {
    return api.post('/Api/Donation', donation);
  },
  onSuccess: () => {
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Thank you for donating!',
      life: 3000,
    });

    queryClient.invalidateQueries({ queryKey: ['meals'] });
    setNameCookie(name.value);
    router.push('/');
  },
  onError: () => {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Unable to donate meal',
      life: 3000,
    });
  }
});

const validateDonationForm = (name: string, selectedMealType: number): boolean => {
  let valid = true;

  if (!name || name.trim() === '') {
    userNameInputErrorText.value = 'Please enter a name';
    valid = false;
  } else if (!/^(\w+\s?){1,5}$/.test(name)) {
    userNameInputErrorText.value = 'Please enter a valid name';
    valid = false;
  } else {
    userNameInputErrorText.value = '';
  }

  if (!selectedMealType || selectedMealType === 0) {
    mealInputErrorText.value = 'Please select a meal';
    valid = false;
  } else {
    mealInputErrorText.value = '';
  }

  return valid;
};

const submitMeal = () => {
  const valid = validateDonationForm(name.value, selectedMealType.value);
  if (!valid) return;

  donationMutation.mutate({
    donorName: name.value,
    mealId: selectedMealType.value
  });
};
</script>

<style scoped>
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

  .full-width {
    width: 100%;
  }

  .error-text {
    text-align: left;
  }
</style>
