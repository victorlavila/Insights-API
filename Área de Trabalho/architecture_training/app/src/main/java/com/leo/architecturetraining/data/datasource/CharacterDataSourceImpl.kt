package com.leo.architecturetraining.data.datasource

import com.leo.architecturetraining.data.model.CharacterResponse
import com.leo.architecturetraining.data.service.RickAndMortyApi

class CharacterDataSourceImpl(
    private val service: RickAndMortyApi
) : CharacterDataSource {
    override suspend fun getAllCharacter(): CharacterResponse {
        return service.getAllCharacters()
    }
}