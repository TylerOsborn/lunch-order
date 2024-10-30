<template>
  <div class="flex">
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
  import api from '../axios/axios';
  import { ApiResult, User } from '../models/models';
  import {
    setUUIDToLocalStorage,
    setNameToLocalStorage,
    getNameFromLocalStorage,
    getUUIDFromLocalStorage,
  } from '../utils/utils';
  import { UnclaimedDonation, Donation } from '../models/models';

  export default {
    name: 'RecipientForm',
    components: {
      Listbox,
      Button,
      InputText,
      Dialog,
    },
    props: {
      availableMeals: {
        type: Array<UnclaimedDonation>,
        required: false,
        default: [],
      },
    },
    data() {
      return {
        availableMeals: [] as UnclaimedDonation[] | null,
        claimedMeal: null as Donation | null,
        selectedDonation: {
          description: '',
          donorName: '',
        } as UnclaimedDonation,
        name: '' as string,
        uuid: '' as string,
        dialogVisible: false as boolean,

        userNameInputErrorText: '',
        mealInputErrorText: '',
      };
    },
    mounted() {
      this.name = getNameFromLocalStorage();
      this.uuid = getUUIDFromLocalStorage();
    },
    methods: {
      handleOkayButton() {
        this.$router.push('/');
      },
      validateDonationClaim(name: string, donation: UnclaimedDonation): boolean {
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
      selectMeal(donation: UnclaimedDonation) {
        const valid = this.validateDonationClaim(this.name, donation);
        if (!valid) {
          return;
        }
        api
          .post('/Api/Donation/Claim', {
            donationId: donation.id,
            user: {
              uuid: this.uuid,
              name: this.name,
            },
          })
          .then((response) => {
            const result: ApiResult<User> = response.data;
            const uuid = result.data.uuid;
            setUUIDToLocalStorage(uuid);

            if (response.status === 200) {
              this.dialogVisible = true;
              return;
            }
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000 });
            this.$emit('reloadMeals');
          })
          .catch((_) => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Unable to claim meal', life: 3000 });
            this.$emit('reloadMeals');
          });

        setNameToLocalStorage(this.name);
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
