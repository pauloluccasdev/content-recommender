"""
Funções utilitárias para cálculo de similaridade entre usuários
Usa Scikit-learn para calcular similaridade de cosseno
"""

import pandas as pd
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity
from typing import List, Tuple

def calculate_user_similarity(
    interactions_matrix: pd.DataFrame,
    user_id: int
) -> pd.Series:
    """
    Calcula similaridade entre um usuário e todos os outros usuários
    usando similaridade de cosseno

    Args:
        interactions_matrix: Matriz de interações usuário-conteúdo
        user_id: ID do usuário para cancelar similaridade

    Returns:
        pd.Series: Similaridade ordenadas (maior para menor)
    """
    
    if user_id not in interactions_matrix.index:
        return pd.Series(dtype=float)

    # Pegar o vetor do usuário
    user_vector = interactions_matrix.loc[user_id].values.reshape(1, -1)
    
    # Calcular similaridade com todos os usuários
    similarities = cosine_similarity(user_vector, interactions_matrix.values)

    # Converter para Series indexada por user_id
    similarity_series = pd.Series(
        similarities[0],
        index=interactions_matrix.index
    )

    # Remover o próprio usuário e ordenar
    similarity_series = similarity_series.drop(user_id)
    similarity_series = similarity_series.sort_values(ascending=False)

    return similarity_series

def get_top_similar_users(
    similarities: pd.Series,
    top_n: int= 5
) -> List[Tuple[int, float]]:
    """
    Retorna os top N usuários mais similares

    Args:
        similarities: Series com similaridades
        top_n: Número de usuários similares a retornar

    Returns:
        List[Tuple[int, float]]: Lista de (user_id, similarity_score)
    """
    return [
        (user_id, float(score))
        for user_id, score in similarities.head(top_n).items()
    ]