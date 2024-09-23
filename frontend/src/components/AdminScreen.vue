<script lang="ts">
import api from "../axios/axios.ts";
import {ApiResult, Meal, MealType, MenuItem} from "../models/models.ts";

import Card from 'primevue/card';
import Listbox from 'primevue/listbox';
import FileUpload from 'primevue/fileupload';
import Divider from 'primevue/divider';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';


export default {
  name: 'AdminScreen',
  components: {
    Card,
    Listbox,
    FileUpload,
    Divider,
    DataTable,
    Column,
  },
  data() {
    return {
      meals: [] as Meal[],
      mealTypes: [] as MealType[],
      weeklyMeals: [] as MenuItem[]
    }
  },
  mounted() {
    this.getMeals();
    this.getMealTypes();
    this.getWeeklyMeals();
  },
  methods: {
    getMeals() {
      api.get(`/Api/Meal?startDate=${this.mondayDate}&endDate=${this.thursdayDate}`)
          .then(response => {
            let result: ApiResult<Meal[]> = response.data;
            this.meals = result.data;
          })
          .catch(error => {
            console.log(error);
          });
    },
    getMealTypes() {
      api.get('/Api/MealType')
          .then(response => {
            let result: ApiResult<MealType[]> = response.data;
            this.mealTypes = result.data;
          })
          .catch(error => {
            console.log(error);
          });
    },
    getWeeklyMeals() {
      api.get(`/Api/Menu?startDate=${this.mondayDate}`)
          .then(response => {
            let result: ApiResult<MenuItem[]> = response.data;
            this.weeklyMeals = result.data;
          })
    },
    zeroPad(num, places) {
      const zero = places - num.toString().length + 1;
      return Array(+(zero > 0 && zero)).join("0") + num;
    },
  },
  computed: {
    mondayDate() {
      const today = new Date();
      const day = today.getDay();
      const diff = today.getDate() - day + (day == 0 ? -6 : 1);
      const monday = new Date(today.setDate(diff))
      const year = this.zeroPad(monday.getFullYear(), 4);
      const month = this.zeroPad(monday.getMonth() + 1, 2);
      const date = this.zeroPad(monday.getDate(), 2);
      return `${year}-${month}-${date}`
    },
    thursdayDate() {
      const today = new Date();
      const day = today.getDay();
      const diff = today.getDate() - day + (day == 0 ? -6 : 4);
      const monday = new Date(today.setDate(diff))
      const year = this.zeroPad(monday.getFullYear(), 4);
      const month = this.zeroPad(monday.getMonth() + 1, 2);
      const date = this.zeroPad(monday.getDate(), 2);
      return `${year}-${month}-${date}`
    },
  }
}
</script>

<template>
  <div class="container">
    <Card class="card">
      <template #title>
        <h2>
          All Meals
        </h2>
      </template>
      <template #content>
        <Listbox filter :options="mealTypes" optionLabel="description"/>
        <Divider/>
        <h3>Upload New Meals</h3>
        <FileUpload mode="basic" name="demo[]" url="/Api/MealType/Upload" accept="csv/*" :maxFileSize="1000000"
                    @upload="uploadMealTypes" :auto="true" chooseLabel="Browse"/>
      </template>
    </Card>
    <Card class="card">
      <template #title>
        <h2>
          This weeks meals
        </h2>
      </template>
      <template #content>
        <DataTable :value="weeklyMeals" scrollable scrollHeight="400px">
          <Column field="date" header="Date"/>
          <Column field="description" header="Description"/>
        </DataTable>
        <Divider/>
        <h3> Upload Weekly Meals</h3>
        <FileUpload mode="basic" name="demo[]" url="/Api/Meal/Upload" accept="csv/*" :maxFileSize="1000000"
                    @upload="uploadMeals" :auto="true" chooseLabel="Browse"/>
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

.card {
  height: 100%;
  width: 50%;
}

#app {
  max-width: 100vw !important;
}
</style>