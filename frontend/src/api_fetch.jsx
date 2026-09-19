


export async  function apiFetch(url, options={}) {
    let accessToken = localStorage.getItem("accessToken");   
    const apiResponse = await fetch(url, {
        ...options,
        headers: {
            ...options.headers,
            'Authorization': `Bearer: ${accessToken}`,
            'Content-Type': "application/json"
        }
    });
    
    if(apiResponse.status != 401) return apiResponse;

    const refreshToken = localStorage.getItem("refreshToken");
     if(!refreshToken) {
        return apiResponse;
    }
    const refreshResponse = await fetch(`http://localhost:8080/api/refresh`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${refreshToken}`,
            'Content-type': 'application/json'
        }
    });

    if(!refreshResponse.ok) {
        localStorage.clear();
        return apiResponse;
    }
    const data =  await refreshResponse.json();
    const newAccessToken = data.token;
    localStorage.setItem("accessToken", newAccessToken);
    const retryResponse = await fetch(url, {
        ...options, 
        headers: {
            ...options.headers,
            'Authorization' : `Bearer ${newAccessToken}`,
            'Content-Type': 'application/json'
        }
    });
    if(retryResponse.status == 401) localStorage.clear();
    return retryResponse;
}