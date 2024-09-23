<template>
  <div class="give-meal-screen">
    <h2>Give a Meal</h2>
    <form class="flex" @submit.prevent="submitMeal">
      <InputText class="full-width" placeholder="Name" id="name" v-model="name"/>
      <Listbox class="full-width" size v-model="selectedMealType" :options="meals" optionValue="description" optionLabel="description" placeholder="Select..." id="meal" required/>
      <Button class="full-width" type="submit">Submit</Button>
    </form>
  </div>
</template>

<script lang="ts">
import Listbox from 'primevue/Listbox';
import FloatLabel from 'primevue/floatlabel';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import {ApiResult, Meal} from "../models/models.ts";
import api from "../axios/axios.ts";
import {getNameFromCookie, setNameCookie} from "../utils/utils.ts";

export default {
  name: 'GiveMealScreen',
  components: {
    Listbox,
    FloatLabel,
    Button,
    InputText
  },
  data() {
    return {
      name: '',
      selectedMealType: '',
      meals: [] as Meal[]
    }
  },
  mounted() {
    this.getMeals();
    this.name = getNameFromCookie();
  },
  methods: {
    getMeals() {
      api.get(`/Api/Meal/Today`)
          .then(response => {
            let result: ApiResult<Meal[]> = response.data;
            this.meals = result.data;
          })
          .catch(error => {
            console.log(error);
          });
    },
    submitMeal() {
      console.log('Meal offered:', {name: this.name, mealType: this.selectedMealType})
      setNameCookie(this.name)
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.flex {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  justify-content: center;
  align-items: center;
}

.full-width {
  width: 100%;
}
</style>