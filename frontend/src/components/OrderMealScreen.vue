<template>
  <div class="order-meal-screen">
    <h2>Order Meals for the Week</h2>

    <div v-if="existingOrder">
      <!-- Show summary if already ordered -->
      <div class="summary-container">
        <h3>Your Meal Order Summary</h3>
        <p class="info-text">You have already submitted your meal order for this week.</p>

        <div class="meal-summary">
          <div class="day-summary">
            <strong>Monday:</strong>
            <span>{{ existingOrder.mondayMeal?.description || 'No meal selected' }}</span>
          </div>
          <div class="day-summary">
            <strong>Tuesday:</strong>
            <span>{{ existingOrder.tuesdayMeal?.description || 'No meal selected' }}</span>
          </div>
          <div class="day-summary">
            <strong>Wednesday:</strong>
            <span>{{ existingOrder.wednesdayMeal?.description || 'No meal selected' }}</span>
          </div>
          <div class="day-summary">
            <strong>Thursday:</strong>
            <span>{{ existingOrder.thursdayMeal?.description || 'No meal selected' }}</span>
          </div>
        </div>

        <Button @click="goHome" class="back-button">Back to Home</Button>
      </div>
    </div>

    <div v-else>
      <!-- Show order form if not yet ordered -->
      <p class="info-text">Select at most one meal per day (Monday-Thursday). Week starting: {{ weekStartDate }}</p>

      <form @submit.prevent="handleSubmit" class="order-form">
        <div v-for="day in days" :key="day.name" class="day-selection">
          <h3>{{ day.name }}</h3>
          <Listbox
            v-model="selections[day.key]"
            :options="getMealsForDay(day.name)"
            optionValue="id"
            optionLabel="description"
            placeholder="No meal selected"
            class="meal-listbox"
          />
        </div>

        <Button
          type="submit"
          :disabled="isSubmitting || isLoadingMeals"
          class="submit-button"
        >
          {{ isSubmitting ? 'Submitting...' : 'Submit Meal Order' }}
        </Button>
      </form>
    </div>

    <!-- Confirmation dialog -->
    <Dialog v-model:visible="showConfirmDialog" header="Confirm Order" :modal="true">
      <p>You have not selected meals for the following day(s):</p>
      <ul>
        <li v-for="day in daysWithoutMeals" :key="day">{{ day }}</li>
      </ul>
      <p>Are you sure you want to proceed?</p>
      <template #footer>
        <Button label="Cancel" @click="showConfirmDialog = false" severity="secondary" />
        <Button label="Yes, Submit" @click="confirmSubmit" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useQuery, useMutation } from '@tanstack/vue-query';
import { useRouter } from 'vue-router';
import type { Meal, ApiResult, MealOrder, MealOrderRequest } from '../models/models';
import Listbox from 'primevue/listbox';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import api from '../axios/axios';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const router = useRouter();

const selections = ref<{
  monday: number | null;
  tuesday: number | null;
  wednesday: number | null;
  thursday: number | null;
}>({
  monday: null,
  tuesday: null,
  wednesday: null,
  thursday: null,
});

const showConfirmDialog = ref(false);
const daysWithoutMeals = ref<string[]>([]);
const isSubmitting = ref(false);
const existingOrder = ref<MealOrder | null>(null);

// Get Monday of current or next week
const getWeekStartDate = (): string => {
  const today = new Date();
  const dayOfWeek = today.getDay();
  const diff = dayOfWeek === 0 ? 1 : (dayOfWeek <= 4 ? 1 - dayOfWeek : 8 - dayOfWeek);
  const monday = new Date(today);
  monday.setDate(today.getDate() + diff);
  return monday.toISOString().split('T')[0];
};

const weekStartDate = ref(getWeekStartDate());

const days = [
  { name: 'Monday', key: 'monday' as const },
  { name: 'Tuesday', key: 'tuesday' as const },
  { name: 'Wednesday', key: 'wednesday' as const },
  { name: 'Thursday', key: 'thursday' as const },
];

// Fetch meals for the week
const { data: mealsResult, isLoading: isLoadingMeals } = useQuery({
  queryKey: ['meals', 'week', weekStartDate.value],
  queryFn: async (): Promise<ApiResult<Meal[]>> => {
    const endDate = new Date(weekStartDate.value);
    endDate.setDate(endDate.getDate() + 3); // Thursday
    const response = await api.get('/Api/Meal', {
      params: {
        startDate: weekStartDate.value,
        endDate: endDate.toISOString().split('T')[0],
      },
    });
    return response.data;
  },
});

// Check for existing order
const checkExistingOrder = async () => {
  try {
    const response = await api.get<ApiResult<MealOrder>>('/Api/MealOrder', {
      params: {
        weekStartDate: weekStartDate.value,
      },
    });
    if (response.data.data) {
      existingOrder.value = response.data.data;
    }
  } catch (error) {
    console.error('Error checking existing order:', error);
  }
};

onMounted(() => {
  checkExistingOrder();
});

const meals = computed(() => mealsResult.value?.data || []);

const getMealsForDay = (dayName: string): Meal[] => {
  const dayDate = new Date(weekStartDate.value);
  const dayIndex = days.findIndex(d => d.name === dayName);
  dayDate.setDate(dayDate.getDate() + dayIndex);
  const dateStr = dayDate.toISOString().split('T')[0];

  return meals.value.filter(meal => meal.date === dateStr);
};

const orderMutation = useMutation({
  mutationFn: async (order: MealOrderRequest) => {
    return api.post('/Api/MealOrder', order);
  },
  onSuccess: () => {
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Your meal order has been submitted!',
      life: 3000,
    });
    checkExistingOrder(); // Refresh to show summary
  },
  onError: (error: any) => {
    const errorMessage = error.response?.data?.error || 'Unable to submit meal order';
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: errorMessage,
      life: 3000,
    });
    isSubmitting.value = false;
  },
});

const handleSubmit = () => {
  // Check if any days don't have meals selected
  const emptyDays: string[] = [];
  days.forEach(day => {
    if (!selections.value[day.key]) {
      emptyDays.push(day.name);
    }
  });

  if (emptyDays.length > 0) {
    daysWithoutMeals.value = emptyDays;
    showConfirmDialog.value = true;
  } else {
    submitOrder();
  }
};

const confirmSubmit = () => {
  showConfirmDialog.value = false;
  submitOrder();
};

const submitOrder = () => {
  isSubmitting.value = true;

  const orderRequest: MealOrderRequest = {
    weekStartDate: weekStartDate.value,
    mondayMealId: selections.value.monday,
    tuesdayMealId: selections.value.tuesday,
    wednesdayMealId: selections.value.wednesday,
    thursdayMealId: selections.value.thursday,
  };

  orderMutation.mutate(orderRequest);
};

const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.order-meal-screen {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.info-text {
  margin-bottom: 2rem;
  color: #666;
}

.order-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.day-selection {
  border: 1px solid #ddd;
  padding: 1rem;
  border-radius: 8px;
}

.day-selection h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
}

.meal-listbox {
  width: 100%;
}

.submit-button {
  margin-top: 1rem;
  padding: 0.75rem;
}

.summary-container {
  text-align: center;
}

.meal-summary {
  margin: 2rem 0;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 8px;
  text-align: left;
}

.day-summary {
  padding: 0.75rem 0;
  border-bottom: 1px solid #ddd;
  display: flex;
  justify-content: space-between;
}

.day-summary:last-child {
  border-bottom: none;
}

.day-summary strong {
  min-width: 120px;
}

.back-button {
  margin-top: 1rem;
}
</style>
