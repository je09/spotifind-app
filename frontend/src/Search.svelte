<script>
    import {
        WindowSetSize,
        EventsOn,
        LogInfo
    } from "../wailsjs/runtime/runtime.js";
    import {
        ReturnResults,
        ReturnProgress
    } from "../wailsjs/go/main/SpotifindApp.js";
    export let results = [];
    let progress = { Done: 0, Total: 100 }; // Initialize progress

    LogInfo("Search results screen loaded");

    WindowSetSize(600, 700);
    EventsOn("rcv:searchResult", (msg) => {
        let r = JSON.parse(msg);
        LogInfo("Results received: " + JSON.stringify(r.Playlist));
        results = [...results, r];
    });
    EventsOn("rcv:progress", (msg) => {
        let p = JSON.parse(msg);
        if (p.Done > p.Total) {
            p.Total = p.Done;
        }
        progress = p;
    });

    ReturnResults();
    ReturnProgress();

</script>

<div class="results-window">
    <div class="results">
        <h3 class="result-text">Results</h3>
        <p class="result-text">Progress: {progress.Done} of {progress.Total}</p>
        <div class="progress-indicator segmented">
            <div class="progress-indicator-bar" style="width: {progress.Done / progress.Total * 100}%"></div>
        </div>
    </div>
    <div class="table-outer">
        <table>
            <thead>
            <tr>
                <th>Name</th>
                <th>Followers</th>
                <th>Styles</th>
                <th>Contacts</th>
                <th>Description</th>
                <th>Market</th>
                <th>Link</th>
            </tr>
            </thead>
            <tbody>
            {#each results as result}
                <tr>
                    <td>{result.Playlist.Name}</td>
                    <td>{result.Playlist.FollowersTotal}</td>
                    <td>{result.Playlist.Styles}</td>
                    <td>{result.Playlist.Contacts}</td>
                    <td>{result.Playlist.Description}</td>
                    <td>{result.Playlist.Region}</td>
                    <td><a href={result.ExternalURLs} target="_blank">Link</a></td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>

<style>
    .result-text {
        font-family: "Pixelated MS Sans Serif", Arial, sans-serif;
    }

    .results {
        height: 180px;
    }

    .results-window {
        overflow: hidden; /* Hide scroll on the main screen */
        display: flex;
        flex-direction: column;
        height: 100vh; /* Full viewport height */
    }

    .table-outer {
        flex-grow: 1; /* Take up remaining space */
        overflow-y: auto; /* Enable vertical scrolling */
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th, td {
        border: 1px solid;
        padding: 5px;
        text-align: left;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 25ch;
    }
</style>