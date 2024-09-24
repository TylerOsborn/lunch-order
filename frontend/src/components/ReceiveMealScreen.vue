<template>
  <div>
    <h2>Receive a Meal</h2>
    <div v-if="availableMeals === null || availableMeals.length === 0">
      <p>There are no meals available at the moment.</p>
    </div>
    <div class="flex" v-else>
      <InputText class="full-width" placeholder="Name" id="name" v-model="name"/>
      <Listbox v-model="selectedDonation" :options="availableMeals" optionLabel="description"/>
      <Button @click="selectMeal(selectedDonation)" class="full-width">
        Select Option
      </Button>
    </div>
    <Dialog :visible="dialogVisible" modal header="Meal Claimed!">
        <p>You have claimed "{{selectedDonation.description}}" from {{selectedDonation.name}}</p>
        <Button @click="handleOkayButton" label="Okay"/>
    </Dialog>
  </div>
</template>

<script lang="ts">
import Listbox from 'primevue/listbox';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import api from "../axios/axios.ts";
import {ApiResult, Donation} from "../models/models.ts";
import {getNameFromCookie} from "../utils/utils.ts";

export default {
  name: 'ReceiveMealScreen',
  components: {
    Listbox,
    Button,
    InputText,
    Dialog
  },
  data() {
    return {
      availableMeals: [] as string[] | null,
      selectedDonation: {
        description: 'Meal',
        name: 'John Doe',
      } as Donation,
      name: '' as string,
      dialogVisible: false as boolean
    }
  },
  mounted() {
    this.getAvailableMeals();
    this.name = getNameFromCookie();
  },
  methods: {
    getAvailableMeals() {
      api.get('/Api/Donation')
          .then(response => {
            let result: ApiResult<Donation[]> = response.data;
            this.availableMeals = result.data;
          })
          .catch(_ => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
          });
    },
    handleOkayButton() {
      this.$router.push('/');
    },
    selectMeal(donation) {
      api.post('/Api/Donation/Claim', {
        donationId: donation.id,
        name: this.name
      })
          .then(response => {
            if (response.status === 200) {
              this.dialogVisible = true;
              return
            }
            this.$toast.add({severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000});
            this.getAvailableMeals();
          })
          .catch(_ => {
            this.$toast.add({severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000});
            this.getAvailableMeals();
          });
    }
  }
}
</script>

<style scoped>
.full-width {
  width: 100%;
}

.flex {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  justify-content: center;
  align-items: center;
}
</style>