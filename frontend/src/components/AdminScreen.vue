<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import api from '../axios/axios.ts';
import { ApiResult, DonationClaimSummary, Meal } from '../models/models.ts';
import { getTodayDate, mondayDate, thursdayDate } from '../utils/utils.ts';
import { useAuth } from '../composables/useAuth.ts';

import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import Divider from 'primevue/divider';
import { useToast } from 'primevue/usetoast';
import AdminLogin from './AdminLogin.vue';

const toast = useToast();
const queryClient = useQueryClient();
const { isAuthenticated, logout } = useAuth();

const newMeals = ref('');

const monday = computed(() => mondayDate());
const thursday = computed(() => thursdayDate());
const today = computed(() => getTodayDate());

const { data: meals = [] } = useQuery({
  queryKey: ['meals', monday.value, thursday.value],
  queryFn: async () => {
    const { data } = await api.get(`/Api/Meal?startDate=${monday.value}&endDate=${thursday.value}`);
    const result: ApiResult<Meal[]> = data;
    return result.data;
  }
});

const { data: claimsSummary = [] } = useQuery({
  queryKey: ['claimsSummary', today.value],
  queryFn: async () => {
    try {
      const { data } = await api.get(`/Api/Stats/Claims/Summary?date=${today.value}&timestamp=${new Date().getTime()}`);
      const result: ApiResult<DonationClaimSummary[]> = data;

      if (result.error) {
        toast.add({ severity: 'error', summary: 'Error', detail: result.error });
        return [];
      }

      return result.data;
    } catch (error) {
      toast.add({ severity: 'error', summary: 'Error', detail: `Error: ${error}` });
      return [];
    }
  }
});

const { mutate: submitMeal } = useMutation({
  mutationFn: async () => {
    return api.post('/Api/Meal/Upload', { csv: newMeals.value });
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['meals'] });
    newMeals.value = '';
  },
  onError: (error) => {
    console.error(error);
    toast.add({ severity: 'error', summary: 'Error', detail: `Error: ${error}` });
  }
});

const handleSubmitMeal = () => {
  submitMeal();
};

const handleLogin = () => {
  // Refresh queries after login
  queryClient.invalidateQueries({ queryKey: ['meals'] });
  queryClient.invalidateQueries({ queryKey: ['claimsSummary'] });
};

const handleLogout = () => {
  logout();
  toast.add({ 
    severity: 'info', 
    summary: 'Logged out', 
    detail: 'You have been logged out successfully',
    life: 3000 
  });
};
</script>

<template>
  <div v-if="!isAuthenticated">
    <AdminLogin @login="handleLogin" />
  </div>
  <div v-else class="container">
    <div class="admin-header">
      <h1>Admin Dashboard</h1>
      <Button 
        label="Logout" 
        icon="pi pi-sign-out" 
        severity="secondary" 
        @click="handleLogout" 
        class="logout-button"
      />
    </div>
    <Card class="card">
      <template #title>
        <h2>Daily Summary</h2>
      </template>
      <template #content>
        <div style="height: inherit">
          <DataTable scrollable scrollHeight="400px" :value="claimsSummary">
            <Column field="claimed" header="Claimed" />
            <Column field="description" header="Description" />
            <Column field="donorName" header="Donor" />
            <Column field="recipientName" header="Recipient" />
          </DataTable>
        </div>
      </template>
    </Card>
    <Card class="card">
      <template #title>
        <h2>This weeks meals</h2>
      </template>
      <template #content>
        <DataTable :value="meals" scrollable scrollHeight="400px">
          <Column field="date" header="Date" />
          <Column field="description" header="Description" />
        </DataTable>
        <Divider />
        <h3>Upload Weekly Meals</h3>
        <form>
          <Textarea rows="10" cols="72" v-model="newMeals" placeholder="Enter weekly meals here" />
          <Button type="submit" class="sub-button" @click.prevent="handleSubmitMeal">Submit</Button>
        </form>
      </template>
    </Card>
  </div>
</template>

<style>
  .admin-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding: 0 1rem;
  }

  .admin-header h1 {
    margin: 0;
  }

  .logout-button {
    margin-left: auto;
  }

  .container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 2rem;
  }

  .container .card {
    height: 100%;
    width: 50%;
  }

  @media (min-width: 768px) {
    .container {
      flex-direction: row;
      justify-content: center;
      align-items: flex-start;
      height: calc(100vh - 4rem);
      width: calc(100vw - 4rem);
    }
  }

  .sub-button {
    margin: 10px;
    width: 100%;
  }

  #app {
    max-width: 100vw !important;
  }
</style>
