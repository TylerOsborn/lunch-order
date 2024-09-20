<template>
  <div class="receive-meal-screen">
    <h2>Receive a Meal</h2>
    <div v-if="availableMeals === null || availableMeals.length === 0">
      <p>There are no meals available at the moment.</p>
    </div>
    <Listbox v-else v-model="selectedMeal" :options="availableMeals" optionLabel="type"/>
    <Button v-if="availableMeals !== null && availableMeals.length > 0" @click="selectMeal(selectedMeal)" class="full-width">
      Select Option
    </Button>
  </div>
</template>

<script>
import Listbox from 'primevue/listbox';
import Button from 'primevue/button';

export default {
  name: 'ReceiveMealScreen',
  components: {
    Listbox,
    Button
  },
  data() {
    return {
      availableMeals: [
        { id: 1, type: 'BAGUETTE - Chicken Coronation Baguette - baguette filled with curried chicken mayo', offeredBy: 'John' },
        { id: 2, type: 'BEEF - Mexican Burrito – Taco spice ground beef, Avo, feta, baby spinach, tomato salsa', offeredBy: 'Jane' },
        { id: 3, type: 'BEEF - Beef Burger & Chips A simple beef patty burger served with potato wedges garnished with lettuce, tomato, and pickles', offeredBy: 'Bob' },
        { id: 3, type: 'BEEF - Korean Crispy Beef Bowl Crispy fried beef strips served with sticky rice, sushi mayo, cucumbers, and a soy reduction', offeredBy: 'Bob' },
      ],
      selectedMeal: null
    }
  },
  methods: {
    selectMeal(meal) {
      alert(`You've selected ${meal.type} offered by ${meal.offeredBy}`)
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.full-width {
  width: 100%;
}

.receive-meal-screen {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  justify-content: center;
  align-items: center;
}
</style>