<template>
  <div>
    <h2>Receive a Meal</h2>
    <div v-if="availableMeals === null || availableMeals.length === 0">
      <p>There are no meals available at the moment.</p>
    </div>
    <div v-else class="flex">
      <div class="flex-left full-width">
        <InputText
          id="name"
          v-model="name"
          :invalid="userNameInputErrorText !== ''"
          class="full-width"
          placeholder="Name"
        />
        <small v-if="userNameInputErrorText !== ''" class="error-text">{{ userNameInputErrorText }}</small>
      </div>
      <div class="flex-left full-width">
        <Listbox
          v-model="selectedDonation"
          :invalid="mealInputErrorText !== ''"
          :options="availableMeals"
          optionLabel="description"
        />
        <small v-if="mealInputErrorText !== ''" class="error-text">{{ mealInputErrorText }}</small>
      </div>
      <Button class="full-width" @click="selectMeal(selectedDonation)"> Select Option</Button>
    </div>
    <Dialog :visible="dialogVisible" header="Meal Claimed!" modal>
      <p>You have claimed "{{ selectedDonation.description }}" from {{ selectedDonation.donorName }}</p>
      <Button label="Okay" @click="handleOkayButton" />
    </Dialog>
  </div>
</template>

<script lang="ts">
  import Listbox from 'primevue/listbox';
  import InputText from 'primevue/inputtext';
  import Button from 'primevue/button';
  import Dialog from 'primevue/dialog';
  import api from '../axios/axios.ts';
  import { ApiResult, Donation } from '../models/models.ts';
  import { getNameFromCookie, setNameCookie } from '../utils/utils.ts';

  export default {
    name: 'ReceiveMealScreen',
    components: {
      Listbox,
      Button,
      InputText,
      Dialog,
    },
    data() {
      return {
        availableMeals: [] as Donation[] | null,
        selectedDonation: {
          description: '',
          donorName: '',
        } as Donation,
        name: '' as string,
        dialogVisible: false as boolean,

        userNameInputErrorText: '',
        mealInputErrorText: '',
      };
    },
    mounted() {
      this.getAvailableMeals();
      this.name = getNameFromCookie();
    },
    methods: {
      getAvailableMeals() {
        api
          .get(`/Api/Donation?timestamp=${new Date().getTime()}`)
          .then((response) => {
            let result: ApiResult<Donation[]> = response.data;
            this.availableMeals = result.data;
          })
          .catch((_) => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
          });
      },
      handleOkayButton() {
        this.$router.push('/');
      },
      validateDonationClaim(name: string, donation: Donation): boolean {
        let valid = true;

        const id = donation?.id;
        const description = donation?.description;

        if (!name || name.trim() === '') {
          this.userNameInputErrorText = 'Please enter a name';
          valid = valid && false;
        } else if (!/^(\w+\s?){1,5}$/.test(name)) {
          this.userNameInputErrorText = 'Please enter a valid name';
          valid = valid && false;
        } else {
          this.userNameInputErrorText = '';
        }

        if (id == null || id <= 0 || description == null || description.trim() === '') {
          this.mealInputErrorText = 'Please select a meal';
          valid = valid && false;
        } else {
          this.mealInputErrorText = '';
        }

        return valid;
      },
      selectMeal(donation: Donation) {
        const valid = this.validateDonationClaim(this.name, donation);
        if (!valid) {
          return;
        }

        api
          .post('/Api/Donation/Claim', {
            donationId: donation.id,
            name: this.name,
          })
          .then((response) => {
            if (response.status === 200) {
              this.dialogVisible = true;
              return;
            }
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000 });
            this.getAvailableMeals();
          })
          .catch((_) => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000 });
            this.getAvailableMeals();
          });
        setNameCookie(this.name);
      },
    },
  };
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

  .flex-left {
    display: flex;
    flex-direction: column;
    justify-content: left;
  }

  .error-text {
    text-align: left;
  }
</style>
