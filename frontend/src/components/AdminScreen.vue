<script lang="ts">
import api from "../axios/axios.ts";
import {ApiResult, Meal } from "../models/models.ts";
import {mondayDate, thursdayDate} from "../utils/utils.ts";

import Card from 'primevue/card';
import Listbox from 'primevue/listbox';
import FileUpload from 'primevue/fileupload';
import Divider from 'primevue/divider';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Textarea from 'primevue/textarea';
import Button from 'primevue/Button';


export default {
  name: 'AdminScreen',
  components: {
    Card,
    Listbox,
    FileUpload,
    Divider,
    DataTable,
    Column,
    Textarea,
    Button
  },
  data() {
    return {
      meals: [] as Meal[],
      newMeals: '' as string
    }
  },
  mounted() {
    this.getMeals();
  },
  methods: {
    submitMeal() {
      api.post('/Api/Meal/Upload', { csv: this.newMeals })
          .then(response => {
            console.log('Meal uploaded:', response.data);
          })
          .then(() => {
            this.getMeals();
            this.newMeals = '';
          })
          .catch(error => {
            console.log(error);
          });
    },
    getMeals() {
      api.get(`/Api/Meal?startDate=${this.monday}&endDate=${this.thursday}`)
          .then(response => {
            let result: ApiResult<Meal[]> = response.data;
            this.meals = result.data;
          })
          .catch(error => {
            console.log(error);
          });
    },
  },
  computed: {
    monday() {
      return mondayDate();
    },
    thursday() {
      return thursdayDate();
    },
  }
}
</script>

<template>
  <div class="container">
    <Card class="card">
      <template #title>
        <h2>
          This weeks meals
        </h2>
      </template>
      <template #content>
        <DataTable :value="meals" scrollable scrollHeight="400px">
          <Column field="date" header="Date"/>
          <Column field="description" header="Description"/>
        </DataTable>
        <Divider/>
        <h3> Upload Weekly Meals</h3>
        <form>
          <Textarea rows="10" cols="72" v-model="newMeals" placeholder="Enter weekly meals here"/>
          <Button type="submit" class="sub-button" @click.prevent="submitMeal">Submit</Button>
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

#app {
  max-width: 100vw !important;
}
</style>