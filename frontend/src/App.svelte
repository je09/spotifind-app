<script>
  import Main from "./Main.svelte";
  import Search from "./Search.svelte";
  import {LogInfo} from "../wailsjs/runtime/runtime.js";

  let ShowMain = true;
  let ShowSearch = false;

  const showResults = () => {
      LogInfo("Showing search results");
      ShowMain = false;
      ShowSearch = true;
  }

  const handleSearch = (event) => {
    ShowMain = false; // Close main screen
    ShowSearch = true; // Show results
  };
</script>

<main>
    {#if ShowMain}
        <Main on:search={handleSearch}/>
    {/if}
    {#if ShowSearch}
        <Search/>
    {/if}
</main>