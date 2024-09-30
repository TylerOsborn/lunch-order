<template>
  <div class="give-meal-screen">
    <h2>Give a Meal</h2>
    <form class="flex" @submit.prevent="submitMeal">
      <InputText class="full-width" placeholder="Name" id="name" v-model="name"/>
      <Listbox class="full-width" size v-model="selectedMealType" :options="meals" optionValue="description"
               optionLabel="description" placeholder="Select..." id="meal" required/>
      <Button class="full-width" type="submit">Submit</Button>
    </form>
  </div>
</template>

<script lang="ts">
import Listbox from 'primevue/listbox';
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
          .catch(_ => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
          });
    },
    submitMeal() {
      api.post("/Api/Donation", {name: this.name, description: this.selectedMealType})
          .then(response => {
            if (response.status === 200) {
              this.$toast.add({severity: 'success', summary: 'Success', detail: 'Thank you for donating!', life: 3000});
              this.$router.push('/')
            } else {
              this.$toast.add({severity: 'error', summary: 'Error', detail: 'Unable to donate meal', life: 3000});
            }
          })
          .catch(_ => {
            this.$toast.add({severity: 'error', summary: 'Error', detail: 'Unable to donate meal', life: 3000});
          });
      setNameCookie(this.name);
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
