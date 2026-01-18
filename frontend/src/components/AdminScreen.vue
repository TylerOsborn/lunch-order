<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import api from '../axios/axios.ts';
import { ApiResult, DonationClaimSummary, Meal } from '../models/models.ts';
import { getSunday, addDays, formatDate } from '../utils/utils.ts';

import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import Divider from 'primevue/divider';
import DatePicker from 'primevue/datepicker';
import Badge from 'primevue/badge';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const queryClient = useQueryClient();

const newMeals = ref('');

const currentDate = ref(new Date());
const summaryDate = ref(new Date());

const startDate = computed(() => {
  const sunday = getSunday(currentDate.value);
  return formatDate(sunday);
});

const endDate = computed(() => {
  const sunday = getSunday(currentDate.value);
  const saturday = addDays(sunday, 6);
  return formatDate(saturday);
});

const formattedSummaryDate = computed(() => formatDate(summaryDate.value));

const { data: meals = [] } = useQuery({
  queryKey: ['meals', startDate, endDate],
  queryFn: async () => {
    const { data } = await api.get(`/Api/Meal?startDate=${startDate.value}&endDate=${endDate.value}`);
    const result: ApiResult<Meal[]> = data;
    return result.data;
  }
});

const prevWeek = () => {
  currentDate.value = addDays(currentDate.value, -7);
};

const nextWeek = () => {
  currentDate.value = addDays(currentDate.value, 7);
};

const resetToToday = () => {
  currentDate.value = new Date();
};

const { data: claimsSummary } = useQuery({
  queryKey: ['claimsSummary', formattedSummaryDate],
  queryFn: async () => {
    try {
      const { data } = await api.get(`/Api/Stats/Claims/Summary?date=${formattedSummaryDate.value}&timestamp=${new Date().getTime()}`);
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

const claimedCount = computed(() => {
  return (claimsSummary.value || []).filter((c: DonationClaimSummary) => c.claimed).length;
});

const { mutate: submitMeal } = useMutation({
  mutationFn: async () => {
    return api.post('/Api/Meal/Upload', { csv: newMeals.value });
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['meals'] });
    newMeals.value = '';
    toast.add({ severity: 'success', summary: 'Success', detail: 'Meals uploaded successfully' });
  },
  onError: (error) => {
    console.error(error);
    toast.add({ severity: 'error', summary: 'Error', detail: `Error: ${error}` });
  }
});

const handleSubmitMeal = () => {
  submitMeal();
};

const printSummary = () => {
  globalThis.print();
};

const placeholderText = "2023-10-27, \"Pizza Day\"\n2023-10-28, \"Taco Tuesday\"";
</script>

<template>
  <div class="container">
    <Card class="card print-section">
      <template #title>
        <div class="header-container">
          <div class="title-group">
            <h2>Daily Summary</h2>
            <Badge :value="claimedCount" severity="success" v-tooltip="'Total Claimed Meals'" />
          </div>
          <div class="controls-group no-print">
            <Button icon="pi pi-print" @click="printSummary" text rounded v-tooltip="'Print Summary'" />
            <DatePicker v-model="summaryDate" dateFormat="yy-mm-dd" showIcon :maxDate="new Date()" class="date-picker-override" />
          </div>
        </div>
      </template>
      <template #content>
        <div style="height: inherit">
          <DataTable scrollable scrollHeight="400px" :value="claimsSummary" class="summary-table">
            <Column field="claimed" header="Claimed">
              <template #body="slotProps">
                <i v-if="slotProps.data.claimed" class="pi pi-check-circle text-green-500"></i>
                <i v-else class="pi pi-times-circle text-gray-400"></i>
              </template>
            </Column>
            <Column field="description" header="Description" />
            <Column field="donorName" header="Donor" />
            <Column field="recipientName" header="Recipient" />
          </DataTable>
        </div>
      </template>
    </Card>
    <Card class="card no-print">
      <template #title>
        <div class="header-container">
          <h2>Weekly Meals</h2>
          <div class="navigation-controls">
            <Button icon="pi pi-chevron-left" @click="prevWeek" text rounded />
            <span class="date-range">{{ startDate }} to {{ endDate }}</span>
            <Button icon="pi pi-chevron-right" @click="nextWeek" text rounded />
            <Button icon="pi pi-calendar" @click="resetToToday" text rounded v-tooltip="'Today'" />
          </div>
        </div>
      </template>
      <template #content>
        <DataTable :value="meals" scrollable scrollHeight="400px">
          <Column field="date" header="Date" />
          <Column field="description" header="Description" />
        </DataTable>
        <Divider />
        <div class="upload-header">
          <h3>Upload Weekly Meals</h3>
          <i class="pi pi-info-circle" v-tooltip="'Format: YYYY-MM-DD, Description\nExample: 2023-10-27, Pizza Day'" style="cursor: help"></i>
        </div>
        <form>
          <Textarea rows="10" cols="72" v-model="newMeals" :placeholder="placeholderText" />
          <Button type="submit" class="sub-button" @click.prevent="handleSubmitMeal">Submit</Button>
        </form>
      </template>
    </Card>
  </div>
</template>

<style>
  .container {
    display: flex;
    flex-direction: row;
    gap: 1rem;
    justify-content: center;
    align-items: center;

    height: calc(100vh - 4rem);
    width: calc(100vw - 4rem);
  }

  .sub-button {
    margin: 10px;
    width: 100%;
  }

  .card {
    height: 100%;
    width: 50%;
  }

  .header-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title-group {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .controls-group {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .navigation-controls {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .date-range {
    font-weight: bold;
    white-space: nowrap;
  }
  
  .date-picker-override {
    width: 150px;
  }

  .upload-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .text-green-500 {
    color: var(--p-green-500);
  }

  .text-gray-400 {
    color: var(--p-gray-400);
  }

  #app {
    max-width: 100vw !important;
  }

  @media print {
    body * {
      visibility: hidden;
    }
    .print-section, .print-section * {
      visibility: visible;
    }
    .print-section {
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      height: auto;
    }
    .no-print {
      display: none !important;
    }
    .card {
      box-shadow: none;
      border: none;
    }
    /* Ensure table prints fully */
    .p-datatable-wrapper {
      overflow: visible !important;
      max-height: none !important;
    }
  }
</style>
