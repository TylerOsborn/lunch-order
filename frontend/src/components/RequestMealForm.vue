<template>
  <div>
    <InputText v-model="name" placeholder="Name" />
    <Listbox v-model="selectedMeals" :options="meals" multiple checkmark optionValue="id" optionLabel="description" />
  </div>
</template>

<script lang="ts">
  import api from '../axios/axios';
  import { ApiResult, Meal } from '../models/models';
  import Listbox from 'primevue/listbox';
  import { getNameFromLocalStorage, getUUIDFromLocalStorage } from '../utils/utils';

  export default {
    name: 'RequestMealForm',
    components: {
      Listbox,
    },
    data() {
      return {
        meals: [] as Meal[],
        name: '' as string,
        uuid: '' as string,
        selectedMeals: [] as number[],
      };
    },
    mounted() {
      this.uuid = getUUIDFromLocalStorage();
      this.name = getNameFromLocalStorage();
      this.getMeals();
    },
    methods: {
      getMeals() {
        api
          .get(`/Api/Meal/Today`)
          .then((response) => {
            let result: ApiResult<Meal[]> = response.data;
            this.meals = result.data;
          })
          .catch((_) => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
          });
      },
    },
  };
</script>

<style scoped></style>
