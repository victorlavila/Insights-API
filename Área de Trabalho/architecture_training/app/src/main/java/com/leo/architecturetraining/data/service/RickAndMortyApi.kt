package com.leo.architecturetraining.data.service

import com.leo.architecturetraining.data.model.CharacterResponse
import retrofit2.http.GET

interface RickAndMortyApi {

    @GET("/character")
    suspend fun getAllCharacters(): CharacterResponse

}