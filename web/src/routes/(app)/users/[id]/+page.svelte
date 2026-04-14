<script lang="ts">
  const { data } = $props();
</script>

<p>Hello World: {data.userData.displayName}</p>

<img src={data.userData.picture.original} alt="" />
