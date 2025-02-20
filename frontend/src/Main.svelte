<script>
    import {
        WindowSetSize,
        EventsOn,
        LogInfo
    } from "../wailsjs/runtime/runtime.js";
    import {Markets} from "../wailsjs/go/main/SpotifindApp.js";
    import {Search} from "../wailsjs/go/main/SpotifindApp.js";
    import {Alert} from "../wailsjs/go/main/SpotifindApp.js";
    import { createEventDispatcher } from 'svelte';

    LogInfo("Main screen loaded");

    const dispatch = createEventDispatcher();

    let searchQuery = ""
    let ignoreQuery = ""
    let marketInfo = ""

    let marketPopular = false
    let marketUnpopular = false
    let allMarkets = []

    async function getMarkets() {
        allMarkets = await Markets()
    }
    getMarkets()

    function PerformSearch() {
        if (marketPopular) {
            marketInfo = 'popular'
        }
        if (marketUnpopular) {
            marketInfo = 'unpopular'
        }

        if (searchQuery === "") {
            Alert("Please enter a search query.")
            return
        }

        if (marketInfo === "") {
            Alert("Please select a market.")
            return
        }

        dispatch("search");
        Search(searchQuery, ignoreQuery, marketInfo)
    }

    // Unselect all radio buttons if a market is selected from the dropdown.
    function UnselectMarketRadioButtons() {
        var marketInfo = document.getElementsByName('marketInfo');
        for (var i = 0; i < marketInfo.length; i++) {
            marketInfo[i].checked = false;
        }
        marketPopular = false
        marketUnpopular = false
    }

    function ClearSpecificMarket() {
        document.getElementById('marketSelector').value = ''
    }

</script>

<div class="main-screen">
    <div class="">
        <label class="block">
            Search Queries:
            <br>
            <input autocomplete="off" bind:value={searchQuery} class="input" id="search" type="text"/>
        </label>

        <label class="block">
            Ignore Queries:
            <br>
            <input autocomplete="off" bind:value={ignoreQuery} class="input" id="search" type="text"/>
        </label>

        <label class="block">
            Market Info
            <br>
            <input on:click={ClearSpecificMarket} type="radio" id="marketInfo" name="marketInfo" bind:value={marketPopular}> Popular
            <input on:click={ClearSpecificMarket} type="radio" id="marketInfo" name="marketInfo" bind:value={marketUnpopular}> Unpopular

            <br>
            <br>
            Specific Market:
            <select on:input={UnselectMarketRadioButtons} bind:value={marketInfo} id="marketSelector">
                {#each allMarkets as market}
                    <option value={market}>{market}</option>
                {/each}
            </select>
        </label>

        <br>
        <button class="button,block" on:click={PerformSearch}>Search</button>
    </div>
</div>