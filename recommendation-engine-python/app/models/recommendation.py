"""
Modelo de recomendação simples baseado em similaridade e popularidade
Usa Pandas + Scikit-learn conforme especificado na tarefa
"""

import pandas as pd
import numpy as np
from typing import List, Dict
from app.services.dataset_service import DatasetService
from app.utils.similarity import calculate_user_similarity, get_top_similar_users

class SimpleRecommendationModel:
    """Modelo de recomendação simples"""

    def __init__(self, dataset_service: DatasetService):
        self.dataset_service = dataset_service
        self.interactions_df = None
        self.contents_df = None
        self.interactions_matrix = None
        self._initialize_model()

    def _initialize_model(self):
        """Inicializa o modelo carregando o dataset (real ou simulado)"""
        self.interactions_df, self.contents_df = \
            self.dataset_service.load_dataset()
        self.interactions_matrix = self.dataset_service.get_interactions_matrix()
    
    def reload_model(self):
        """Recarrega o modelo com dados atualizados (útil após novas interações)"""
        self.interactions_df, self.contents_df = \
            self.dataset_service.reload_dataset()
        self.interactions_matrix = self.dataset_service.get_interactions_matrix()

    def recommend_by_similarity(
        self,
        user_id: int,
        top_n: int = 10
    ) -> List[Dict]:
        """
        Recomenda conteúdos baseado em similaridade de usuários
        (Collaborative Filtering - filtragem colaborativa)
        
        Algoritmo:
        1. Encontra usuários similares ao usuário alvo
        2. Busca conteúdos que usuários similares gostaram
        3. Calcula score baseado na similaridade e ratings
        4. Retorna top N conteúdos
        
        Args:
            user_id: ID do usuário
            top_n: Número de recomendações
            
        Returns:
            List[Dict]: Lista de recomendações com content_id, score, title
        """
        if user_id not in self.interactions_matrix.index:
            return []
        
        # Calcular similaridade com outros usuários
        similarities = calculate_user_similarity(
            self.interactions_matrix,
            user_id
        )
        
        # Pegar top 10 usuários similares
        similar_users = get_top_similar_users(similarities, top_n=10)
        
        # Conteúdos já visualizados pelo usuário
        user_interactions = set(
            self.interactions_df[
                self.interactions_df['user_id'] == user_id
            ]['content_id'].values
        )
        
        # Calcular scores de recomendação
        recommendation_scores = {}
        
        for similar_user_id, similarity_score in similar_users:
            # Interações do usuário similar
            similar_user_interactions = self.interactions_df[
                self.interactions_df['user_id'] == similar_user_id
            ]
            
            for _, interaction in similar_user_interactions.iterrows():
                content_id = interaction['content_id']
                rating = interaction['rating']
                
                # Pular conteúdos já visualizados
                if content_id in user_interactions:
                    continue
                
                # Calcular score ponderado (similaridade * rating)
                weighted_score = similarity_score * rating
                
                if content_id not in recommendation_scores:
                    recommendation_scores[content_id] = 0.0
                
                recommendation_scores[content_id] += weighted_score
        
        # Normalizar scores (0-1)
        if recommendation_scores:
            max_score = max(recommendation_scores.values())
            if max_score > 0:
                recommendation_scores = {
                    k: v / max_score
                    for k, v in recommendation_scores.items()
                }
        
        # Ordenar e pegar top N
        sorted_recommendations = sorted(
            recommendation_scores.items(),
            key=lambda x: x[1],
            reverse=True
        )[:top_n]
        
        # Formatar resultado
        recommendations = []
        for content_id, score in sorted_recommendations:
            content_info = self.contents_df[
                self.contents_df['content_id'] == content_id
            ].iloc[0]
            
            recommendations.append({
                'content_id': int(content_id),
                'score': float(score),
                'title': content_info['title']
            })
        
        return recommendations
    
    def recommend_by_popularity(
        self,
        user_id: int,
        top_n: int = 10
    ) -> List[Dict]:
        """
        Recomenda conteúdos baseado em popularidade
        (Baseado em ratings e número de interações)
        
        Algoritmo:
        1. Calcula score de popularidade para cada conteúdo
        2. Score = média de ratings * log(número de interações)
        3. Remove conteúdos já visualizados pelo usuário
        4. Retorna top N conteúdos mais populares
        
        Args:
            user_id: ID do usuário
            top_n: Número de recomendações
            
        Returns:
            List[Dict]: Lista de recomendações com content_id, score, title
        """
        # Calcular popularidade de cada conteúdo
        content_stats = self.interactions_df.groupby('content_id').agg({
            'rating': ['mean', 'count']
        })
        content_stats.columns = ['avg_rating', 'interaction_count']
        
        # Score de popularidade = média de rating * log(contagem)
        # Usa log para evitar que conteúdos com muitas interações dominem
        content_stats['popularity_score'] = (
            content_stats['avg_rating'] * 
            np.log1p(content_stats['interaction_count']) / 5.0  # Normalizar
        )
        
        # Conteúdos já visualizados pelo usuário
        user_interactions = set(
            self.interactions_df[
                self.interactions_df['user_id'] == user_id
            ]['content_id'].values
        )
        
        # Remover conteúdos já visualizados e ordenar
        available_contents = content_stats[
            ~content_stats.index.isin(user_interactions)
        ]
        
        top_contents = available_contents.nlargest(top_n, 'popularity_score')
        
        # Normalizar scores (0-1)
        max_score = top_contents['popularity_score'].max()
        if max_score > 0:
            top_contents['popularity_score'] = \
                top_contents['popularity_score'] / max_score
        
        # Formatar resultado
        recommendations = []
        for content_id, row in top_contents.iterrows():
            content_info = self.contents_df[
                self.contents_df['content_id'] == content_id
            ].iloc[0]
            
            recommendations.append({
                'content_id': int(content_id),
                'score': float(row['popularity_score']),
                'title': content_info['title']
            })
        
        return recommendations
    
    def recommend(
        self,
        user_id: int,
        top_n: int = 10,
        method: str = "similarity"
    ) -> List[Dict]:
        """
        Método principal de recomendação
        
        Args:
            user_id: ID do usuário
            top_n: Número de recomendações
            method: 'similarity' ou 'popularity'
            
        Returns:
            List[Dict]: Recomendações formatadas
        """
        if method == "popularity":
            return self.recommend_by_popularity(user_id, top_n)
        else:  # similarity (padrão)
            return self.recommend_by_similarity(user_id, top_n)       

        