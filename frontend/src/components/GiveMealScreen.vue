<template>
  <div class="give-meal-screen">
    <h2>Give a Meal</h2>
    <form class="flex" @submit.prevent="submitMeal">
      <InputText placeholder="Name" id="name" v-model="name"/>
      <Select class="full-width" size v-model="selectedMealType" :options="mealTypes" placeholder="Select..." id="mealType" required/>
      <Button class="full-width" type="submit">Submit</Button>
    </form>
  </div>
</template>

<script lang="ts">
import Select from 'primevue/select';
import FloatLabel from 'primevue/floatlabel';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';

export default {
  name: 'GiveMealScreen',
  components: {
    Select,
    FloatLabel,
    Button,
    InputText
  },
  data() {
    return {
      name: '',
      selectedMealType: '',
      mealTypes: ['BEEF', 'CHICKEN', 'VEGAN']
    }
  },
  methods: {
    submitMeal() {
      console.log('Meal offered:', {name: this.name, mealType: this.selectedMealType})
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