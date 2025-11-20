"""
Serviço para gerar dataset simulado de interações usuário-conteúdo
Usa Pandas para criar e manipular os dados
"""

import pandas as pd
import numpy as np
from typing import Tuple

class DatasetService:
    """Serviço para gerenciar datasets simulados"""

    def __init__(self, n_users: int = 100, n_contents: int = 50):
        self.n_users = n_users
        self.n_contents = n_contents
        self.interactions_df = None
        self.contents_df = None

    def generate_simulated_dataset(self) -> Tuple[pd.DataFrame, pd.DataFrame]:
        """
        Gera um dataset simulado de interações usuário-conteúdo

        Returns:
            Tuple[pd.DataFrame, pd.DataFrame]:
                - DataFrame de interações (user_id, content_id, rating)
                - DataFrame de conteúdos (content_id, title)
        """
        np.random.seed(42) # Para reprodutibilidade

        # Criar DataFrame de conteúdos
        self.contents_df = pd.DataFrame({
            'content_id': range(1, self.n_contents + 1),
            'title': [f'Conteúdo {i}' for i in range(1, self.n_contents + 1)]
        })

        # Criar interações simulares (cada usuário interage com ~30% dos conteúdos)
        interactions = []
        n_interactions = int(self.n_users * self.n_contents * 0.3) 

        for _ in range(n_interactions):
            user_id = np.random.randint(1, self.n_users + 1)
            content_id = np.random.randint(1, self.n_contents + 1)
            # Rating de 1 a 5 (simulando likes/visualizações)
            rating = np.random.choice([1, 2, 3, 4, 5], p=[0.1, 0.1, 0.2, 0.3, 0.3])
            interactions.append({
                'user_id': user_id,
                'content_id': content_id,
                'rating': rating
            })

        self.interactions_df = pd.DataFrame(interactions)
        # Remover duplicatas (manter apenas a última interação)
        self.interactions_df = self.interactions_df.drop_duplicates(
            subset=['user_id', 'content_id'], 
            keep='last'
        )

        return self.interactions_df, self.contents_df
    
    def get_interactions_matrix(self) -> pd.DataFrame:
        """
        Converte interações em uma matriz usuário-conteúdo (pivot table)

        Returns:
            pd.DataFrame: Matriz com usuários nas linhas e conteúdos nas colunas
        """

        if self.interactions_df is None:
            self.generate_simulated_dataset()

        # Criar matriz de interações (pivot table)
        matrix = self.interactions_df.pivot_table(
            index='user_id',
            columns='content_id',
            values='rating',
            aggfunc='mean',
            fill_value=0
        )

        return matrix
        