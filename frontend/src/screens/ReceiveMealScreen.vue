<template>
  <div>
    <h2>Receive a Meal</h2>
    <div v-if="claimedMeal !== null">
      <p>You have claimed "{{ claimedMeal?.meal.description }}" from {{ claimedMeal?.donor?.name }}</p>
    </div>
    <div v-else-if="availableMeals === null || availableMeals.length === 0">
      <p>
        There are no meals available right now but you can let us know what meals you would like and we will update this
        page when a meal is found!
      </p>
      <RequestMealForm />
    </div>
    <RecipientForm @reloadMeals="getAvailableMeals" v-else />
  </div>
</template>

<script lang="ts">
  import RecipientForm from '../components/RecipientForm.vue';
  import RequestMealForm from '../components/RequestMealForm.vue';
  import api from '../axios/axios.ts';
  import { ApiResult, Donation } from '../models/models.ts';
  import { getTodayDate, getUUIDFromLocalStorage } from '../utils/utils.ts';
  import { UnclaimedDonation } from '../models/models.ts';

  export default {
    name: 'ReceiveMealScreen',
    components: {
      RecipientForm,
      RequestMealForm,
    },
    data() {
      return {
        claimedMeal: null as Donation | null,
        uuid: '' as string,
        availableMeals: [] as UnclaimedDonation[] | null,
      };
    },
    beforeMount() {
      this.uuid = getUUIDFromLocalStorage();
      this.checkForMealClaim();
      this.getAvailableMeals();
    },
    methods: {
      checkForMealClaim() {
        api.get(`/Api/Donation?recipientUUID=${this.uuid}&date=${getTodayDate()}`).then((response) => {
          let result: ApiResult<Donation> = response.data;
          this.claimedMeal = result.data;
        });
      },
      getAvailableMeals() {
        api
          .get(`/Api/Donation/Unclaimed?timestamp=${new Date().getTime()}`)
          .then((response) => {
            let result: ApiResult<UnclaimedDonation[]> = response.data;
            this.availableMeals = result.data;
          })
          .catch((_) => {
            this.$toast.add({ severity: 'error', summary: 'Error', detail: 'Error loading meal options', life: 3000 });
          });
      },
    },
  };
</script>

<style scoped></style>
