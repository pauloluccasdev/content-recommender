"""
Serviço para gerar dataset simulado ou buscar dados reais do banco
Usa Pandas para criar e manipular os dados
"""

import pandas as pd
import numpy as np
from typing import Tuple, Optional
import logging
from app.core.config import settings
from app.services.database_service import database_service

logger = logging.getLogger(__name__)

class DatasetService:
    """Serviço para gerenciar datasets (reais ou simulados)"""

    def __init__(self, n_users: int = 100, n_contents: int = 50):
        self.n_users = n_users
        self.n_contents = n_contents
        self.interactions_df = None
        self.contents_df = None
        self.use_real_data = settings.data_mode == "real"

    def load_dataset(self) -> Tuple[pd.DataFrame, pd.DataFrame]:
        """
        Carrega dataset - tenta usar dados reais primeiro, fallback para simulados
        
        Returns:
            Tuple[pd.DataFrame, pd.DataFrame]:
                - DataFrame de interações (user_id, content_id, rating)
                - DataFrame de conteúdos (content_id, title)
        """
        if self.use_real_data:
            # Tentar carregar dados reais do banco
            real_data = self._load_real_dataset()
            if real_data[0] is not None and real_data[1] is not None:
                self.interactions_df, self.contents_df = real_data
                logger.info("Usando dados reais do banco de dados")
                return self.interactions_df, self.contents_df
            else:
                logger.warning("Não foi possível carregar dados reais, usando dados simulados")
        
        # Fallback para dados simulados
        return self.generate_simulated_dataset()
    
    def _load_real_dataset(self) -> Tuple[Optional[pd.DataFrame], Optional[pd.DataFrame]]:
        """
        Carrega dados reais do banco de dados
        
        Returns:
            Tuple[Optional[pd.DataFrame], Optional[pd.DataFrame]]:
                - DataFrame de interações ou None
                - DataFrame de conteúdos ou None
        """
        try:
            interactions_df = database_service.fetch_interactions()
            contents_df = database_service.fetch_contents()
            
            if interactions_df is None or contents_df is None:
                return None, None
            
            # Garantir que as colunas necessárias existem
            if 'rating' not in interactions_df.columns:
                interactions_df['rating'] = 3.0
            
            # Remover duplicatas (manter última interação)
            interactions_df = interactions_df.drop_duplicates(
                subset=['user_id', 'content_id'],
                keep='last'
            )
            
            # Garantir que contents_df tem title
            if 'title' not in contents_df.columns or contents_df['title'].isna().all():
                if 'description' in contents_df.columns:
                    contents_df['title'] = contents_df['description'].fillna('Conteúdo')
                elif 'content_type' in contents_df.columns:
                    contents_df['title'] = contents_df['content_type'].fillna('Conteúdo')
                else:
                    contents_df['title'] = 'Conteúdo ' + contents_df['content_id'].astype(str)
            
            return interactions_df, contents_df
            
        except Exception as e:
            logger.error(f"Erro ao carregar dados reais: {e}")
            return None, None

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
    
    def reload_dataset(self) -> Tuple[pd.DataFrame, pd.DataFrame]:
        """
        Recarrega o dataset (útil quando há novas interações)
        
        Returns:
            Tuple[pd.DataFrame, pd.DataFrame]: Dataset atualizado
        """
        return self.load_dataset()
        